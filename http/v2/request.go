package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type requestMiddlewareHandler = func(*Client, *http.Request) *Client
type responseMiddlewareHandler = func(*Client, *http.Request, *Response) *Response

type Client struct {
	http.Client
	safe        bool // 链式操作安全
	reqMdls     []requestMiddlewareHandler
	resMdls     []responseMiddlewareHandler
	headers     map[string]interface{}
	queryArgs   []interface{}
	cookies     map[string]string
	baseURL     string
	contentType string
	proxy       string
	timeout     time.Duration
	retryCount int
	retryInterval time.Duration
	browserMode bool
}

var defaultHttpClient = http.Client{
	Transport: &http.Transport{
		// No validation for https certification of the server in default.
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableKeepAlives: true,
	},
}

// 创建 HTTP Client 实例
func New() *Client {
	return &Client{
		Client: defaultHttpClient,
		safe: true,
		headers: map[string]interface{}{},
		cookies: map[string]string{},
		contentType: "json",
		timeout: 30 * time.Second,
		retryInterval: time.Second,
	}
}

// Safe 链式安全模式，启用后，链式调用时对所有配置的修改都是逐个叠加，并且不会影响链式前配置，默认开启
func (c *Client) Safe(s ...bool) *Client {
	if len(s) == 0 {
		c.safe = true
	} else {
		c.safe = s[0]
	}
	return c
}

func (c *Client) Clone() *Client {
	//ncli := New()
	ncli := &Client{}
	*ncli = *c
	//query
	if l := len(c.queryArgs); l > 0 {
		ncli.queryArgs = make([]interface{}, l)
		copy(ncli.queryArgs, c.queryArgs)
	}

	// headers
	if n := len(c.headers); n > 0 {
		ncli.headers = make(map[string]interface{})
		for k, v := range c.headers {
			ncli.headers[k] = v
		}
	}

	// cookies
	if n := len(c.cookies); n > 0 {
		ncli.cookies = make(map[string]string)
		for k, v := range c.cookies {
			ncli.cookies[k] = v
		}
	}

	// mdls
	if n := len(c.reqMdls); n > 0 {
		ncli.reqMdls = make([]requestMiddlewareHandler, n)
		copy(ncli.reqMdls, c.reqMdls)
	}

	if n := len(c.resMdls); n > 0 {
		ncli.resMdls = make([]responseMiddlewareHandler, n)
		copy(ncli.resMdls, c.resMdls)
	}

	return ncli
}

func (c *Client) getSetting() *Client {
	if c.safe {
		return c.Clone()
	} else {
		return c
	}
}

func (c *Client) setProxy(proxyURL string) {
	if strings.TrimSpace(proxyURL) == "" {
		return
	}
	_proxy, err := url.Parse(proxyURL)
	if err != nil {
		return
	}
	if _proxy.Scheme == "http" {
		if _, ok := c.Transport.(*http.Transport); ok {
			c.Transport.(*http.Transport).Proxy = http.ProxyURL(_proxy)
		}
	} else {
		var auth = &proxy.Auth{}
		user := _proxy.User.Username()

		if user != "" {
			auth.User = user
			password, hasPassword := _proxy.User.Password()
			if hasPassword && password != "" {
				auth.Password = password
			}
		} else {
			auth = nil
		}
		// refer to the source code, error is always nil
		dialer, err := proxy.SOCKS5(
			"tcp",
			_proxy.Host,
			auth,
			&net.Dialer{
				Timeout:   c.Client.Timeout,
				KeepAlive: c.Client.Timeout,
			},
		)
		if err != nil {
			return
		}
		if _, ok := c.Transport.(*http.Transport); ok {
			c.Transport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return dialer.Dial(network, addr)
			}
		}
	}
}

// DoRequest 发起请求， data 参数在 GET, HEAD 请求时被解析为查询参数，在其他方法请求时根据 Content-Type 解析为其他数据类型，如 json, url-form-encoding, xml 等等
func (c *Client) DoRequest(method, u string, data ...interface{}) (resp *Response, err error) {
	cli := c.getSetting()

	method = strings.ToUpper(method)
	u = c.baseURL + strings.Trim(u, " ")

	// header
	headers := http.Header{}
	for k, v := range cli.headers {
		var val string
		switch v.(type) {
		case string:
			val = v.(string)
		case func() string:
			h := v.(func() string)
			val = h()
		case func(*Client) string:
			h := v.(func(*Client) string)
			val = h(cli)
		default:
			continue
		}
		headers.Add(k, val)
	}

	contentType := headers.Get(headerContentTypeKey)
	if contentType == "" {
		contentType = Types[TypeJSON]
		headers.Set(headerContentTypeKey, contentType)
	}

	param := toUrlEncoding(cli.queryArgs...) // url 查询参数

	body := ""  // body 数据
	var doData interface{}
	if len(data) > 0 {
		doData = data[0]
	}

	// 解析第三个参数
	if doData != nil {
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			// 这些请求是没有 body 的
			switch doData.(type) {
			case string:
				param = doData.(string)
			case []byte:
				param = string(doData.([]byte))
			default:
				param = toUrlEncoding(doData)
			}
		} else {
			// 有 body 的酌情序列化
			if strings.Contains(contentType, "/json") {
				switch doData.(type) {
				case string, []byte: // json 字符串
					body = toString(doData)
				default:
					if b, err := json.Marshal(doData); err != nil {
						return nil, err
					} else {
						body = string(b)
					}
				}
			} else if strings.Contains(contentType, "/xml") {
				switch doData.(type) {
				case string, []byte: // xml 字符串
					body = toString(doData)
				default:
					if b, err := xml.Marshal(doData); err != nil {
						return nil, err
					} else {
						body = string(b)
					}
				}
			} else {
				body = toUrlEncoding(data...)
			}
		}
	}
	if param != "" {
		if strings.Contains(u, "?") {
			u += "&" + param
		} else {
			u += "?" + param
		}
	}
	var httpReq *http.Request

	if method == "GET" {
		if httpReq, err = http.NewRequest(method, u, bytes.NewBuffer(nil)); err != nil {
			return nil, err
		} else {
			httpReq.Header = headers
		}
	} else if strings.Contains(body, "@file:") {
		// File uploading request.
		buffer := new(bytes.Buffer)
		writer := multipart.NewWriter(buffer)
		for _, item := range strings.Split(param, "&") {
			array := strings.Split(item, "=")
			if len(array[1]) > 6 && strings.Compare(array[1][0:6], "@file:") == 0 {
				path := array[1][6:]
				if !fileExist(path) {
					return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
				}
				if file, err := writer.CreateFormFile(array[0], filepath.Base(path)); err == nil {
					if f, err := os.Open(path); err == nil {
						if _, err = io.Copy(file, f); err != nil {
							f.Close()
							return nil, err
						}
						f.Close()
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			} else {
				if err = writer.WriteField(array[0], array[1]); err != nil {
					return nil, err
				}
			}
		}
		// Close finishes the multipart message and writes the trailing
		// boundary end line to the output.
		if err = writer.Close(); err != nil {
			return nil, err
		}

		if httpReq, err = http.NewRequest(method, u, buffer); err != nil {
			return nil, err
		} else {
			httpReq.Header = headers
			httpReq.Header.Set("Content-Type", writer.FormDataContentType())
		}
	} else {
		// Normal request.
		paramBytes := []byte(body)
		if httpReq, err = http.NewRequest(method, u, bytes.NewReader(paramBytes)); err != nil {
			return nil, err
		} else {
			httpReq.Header = headers
			if len(paramBytes) > 0 {
				if (paramBytes[0] == '[' || paramBytes[0] == '{') && json.Valid(paramBytes) {
					// Auto detecting and setting the post content format: JSON.
					httpReq.Header.Set("Content-Type", "application/json")
				} else if ok, err := regexp.MatchString(`^[\w\[\]]+=.+`, param); ok && err == nil {
					// If the parameters passed like "name=value", it then uses form type.
					httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				}
			}
		}
	}

	// It's necessary set the cli.Host if you want to custom the host value of the request.
	// It uses the "Host" value from header if it's not set in the request.
	if host := httpReq.Header.Get("Host"); host != "" && httpReq.Host == "" {
		httpReq.Host = host
	}
	// Custom Cookie.
	if len(cli.cookies) > 0 {
		headerCookie := httpReq.Header.Get("Cookie")
		for k, v := range cli.cookies {
			if len(headerCookie) > 0 {
				headerCookie += ";"
			}
			headerCookie += k + "=" + v
		}
		if len(headerCookie) > 0 {
			httpReq.Header.Set("Cookie", headerCookie)
		}
	}

	// HTTP basic authentication.
	//if len(c.authUser) > 0 {
	//	httpReq.SetBasicAuth(c.authUser, c.authPass)
	//}

	resp = newResponse(cli, httpReq)
	// The request body can be reused for dumping
	// raw HTTP request-response procedure.
	//reqBodyContent, _ := ioutil.ReadAll(httpReq.Body)
	//resp.requestBody = reqBodyContent
	//httpReq.Body =  (reqBodyContent)

	if cli.proxy != "" {
		cli.setProxy(cli.proxy)
	}
	cli.Client.Timeout = cli.timeout
	// call middleware
	for _, h := range cli.reqMdls {
		cli = h(cli, httpReq)
	}

	// call res mdl
	defer func() {
		for _, h := range cli.resMdls {
			resp = h(cli, httpReq, resp)
		}
		if err == nil {
			err = resp.GetError()
		}
	}()

	for {
		if resp.Response, err = cli.Do(httpReq); err != nil {
			// The response might not be nil when err != nil.
			if resp.Response != nil {
				resp.Response.Body.Close()
			}
			if cli.retryCount > 0 {
				cli.retryCount--
				time.Sleep(cli.retryInterval)
			} else {
				resp.Err.Add(err)
				return resp, err
			}
		} else {
			break
		}
	}

	// Auto saving cookie content.
	if cli.browserMode {
		now := time.Now()
		for _, v := range resp.Response.Cookies() {
			if !v.Expires.IsZero() && v.Expires.UnixNano() < now.UnixNano() {
				delete(cli.cookies, v.Name)
			} else {
				cli.cookies[v.Name] = v.Value
			}
		}
	}
	resp.ready()
	err = resp.GetError()
	return resp, err
}
