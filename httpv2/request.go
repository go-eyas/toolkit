package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
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

type requestMiddlewareHandler = func(*Request, *http.Request) *Request
type responseMiddlewareHandler = func(*Request, *http.Request, *Response) *Response

type Request struct {
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

func New() *Request {
	return &Request{
		Client: http.Client{
			Transport: &http.Transport{
				// No validation for https certification of the server in default.
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DisableKeepAlives: true,
			},
		},
		safe: true,
		headers: map[string]interface{}{},
		cookies: map[string]string{},
		contentType: "json",
		timeout: 30 * time.Second,
		retryInterval: time.Second,
	}
}

func (r *Request) Safe(s ...bool) *Request {
	if len(s) == 0 {
		r.safe = true
	} else {
		r.safe = s[0]
	}
	return r
}

func (r *Request) Clone() *Request {
	nr := New()
	*nr = *r
	//query
	if l := len(r.queryArgs); l > 0 {
		nr.queryArgs = make([]interface{}, l)
		copy(nr.queryArgs, r.queryArgs)
	}

	// headers
	if n := len(r.headers); n > 0 {
		nr.headers = make(map[string]interface{})
		for k, v := range r.headers {
			nr.headers[k] = v
		}
	}

	// cookies
	if n := len(r.cookies); n > 0 {
		nr.cookies = make(map[string]string)
		for k, v := range r.cookies {
			nr.cookies[k] = v
		}
	}

	// mdls
	if n := len(r.reqMdls); n > 0 {
		nr.reqMdls = make([]requestMiddlewareHandler, n)
		copy(nr.reqMdls, r.reqMdls)
	}

	if n := len(r.resMdls); n > 0 {
		nr.resMdls = make([]responseMiddlewareHandler, n)
		copy(nr.resMdls, r.resMdls)
	}

	return nr
}

func (r *Request) getSetting() *Request {
	if r.safe {
		return r.Clone()
	} else {
		return r
	}
}

func (r *Request) setProxy(proxyURL string) {
	if strings.TrimSpace(proxyURL) == "" {
		return
	}
	_proxy, err := url.Parse(proxyURL)
	if err != nil {
		return
	}
	if _proxy.Scheme == "http" {
		if _, ok := r.Transport.(*http.Transport); ok {
			r.Transport.(*http.Transport).Proxy = http.ProxyURL(_proxy)
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
				Timeout:   r.Client.Timeout,
				KeepAlive: r.Client.Timeout,
			},
		)
		if err != nil {
			return
		}
		if _, ok := r.Transport.(*http.Transport); ok {
			r.Transport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return dialer.Dial(network, addr)
			}
		}
	}
}

func (r *Request) DoRequest(method, u string, data ...interface{}) (resp *Response, err error) {
	req := r.getSetting()
	method = strings.ToUpper(method)
	u = r.baseURL + strings.Trim(u, " ")

	// header
	headers := http.Header{}
	for k, v := range req.headers {
		var val string
		switch v.(type) {
		case string:
			val = v.(string)
		case func() string:
			h := v.(func() string)
			val = h()
		case func(*Request) string:
			h := v.(func(*Request) string)
			val = h(req)
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

	param := "" // url 查询参数
	for _, q := range req.queryArgs {
		param += toUrlEncoding(q)
	}

	body := ""  // body 数据
	var doData interface{}
	if len(data) > 0 {
		doData = data[0]
	}

	// 解析第三个参数
	if doData != nil {
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			switch doData.(type) {
			case string:
				param = doData.(string)
			case []byte:
				param = string(doData.([]byte))
			default:
				param = toUrlEncoding(doData)
			}
		} else {
			switch contentType {
			case "application/json":
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
			case "application/xml": // 这年头谁还用xml啊
				switch doData.(type) {
				case string, []byte: // xml 字符串
					body = toString(doData)
				default:
					//if b, err := gparser.VarToXml(data[0]); err != nil {
					//	return nil, err
					//} else {
					//	param = gconv.UnsafeBytesToStr(b)
					//}
				}

			default: // form 那些
				body = toUrlEncoding(doData)
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
	}else if strings.Contains(body, "@file:") {
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

	// It's necessary set the req.Host if you want to custom the host value of the request.
	// It uses the "Host" value from header if it's not set in the request.
	if host := httpReq.Header.Get("Host"); host != "" && httpReq.Host == "" {
		httpReq.Host = host
	}
	// Custom Cookie.
	if len(req.cookies) > 0 {
		headerCookie := httpReq.Header.Get("Cookie")
		for k, v := range req.cookies {
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

	resp = newResponse(req, httpReq)
	// The request body can be reused for dumping
	// raw HTTP request-response procedure.
	//reqBodyContent, _ := ioutil.ReadAll(httpReq.Body)
	//resp.requestBody = reqBodyContent
	//httpReq.Body =  (reqBodyContent)

	if req.proxy != "" {
		req.setProxy(req.proxy)
	}
	req.Client.Timeout = req.timeout
	// call middleware
	for _, h := range req.reqMdls {
		req = h(req, httpReq)
	}

	// call res mdl
	defer func() {
		for _, h := range req.resMdls {
			resp = h(req, httpReq, resp)
		}
	}()

	for {
		if resp.HttpResponse, err = req.Do(httpReq); err != nil {
			// The response might not be nil when err != nil.
			if resp.HttpResponse != nil {
				resp.HttpResponse.Body.Close()
			}
			if req.retryCount > 0 {
				req.retryCount--
				time.Sleep(req.retryInterval)
			} else {
				resp.Err.Add(err)
				return resp, err
			}
		} else {
			break
		}
	}

	// Auto saving cookie content.
	if req.browserMode {
		now := time.Now()
		for _, v := range resp.HttpResponse.Cookies() {
			if !v.Expires.IsZero() && v.Expires.UnixNano() < now.UnixNano() {
				delete(req.cookies, v.Name)
			} else {
				req.cookies[v.Name] = v.Value
			}
		}
	}
	resp.ready()
	err = resp.GetError()
	return resp, err
}
