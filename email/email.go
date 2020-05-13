package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"sync"

	"github.com/go-eyas/toolkit/util"
	"github.com/jordan-wright/email"
)

type TPL struct {
	Subject string
	From    string
	To      []string
	Bcc     []string
	Cc      []string
	Text    string
	HTML    string
}

type Config struct {
	Name     string
	Account  string
	Password string
	Host     string
	Port     string
	Secure   bool
	TPL      map[string]*TPL
}

type Email struct {
	mailHostAddr string
	SmtpAuth     smtp.Auth
	conf         *Config
	TLS          *tls.Config
	TPLs         map[string]*TPL
	cacheName    string
	sendMu       sync.Mutex
}

func New(conf *Config) *Email {
	auth := smtp.PlainAuth(conf.Name, conf.Account, conf.Password, conf.Host)
	email := &Email{
		mailHostAddr: conf.Host + ":" + conf.Port,
		SmtpAuth:     auth,
		conf:         conf,
		cacheName:    util.RandomStr(6),
	}
	if conf.Secure {
		email.TLS = &tls.Config{
			ServerName:         conf.Host,
			InsecureSkipVerify: true,
		}
	}

	if conf.TPL == nil {
		email.TPLs = make(map[string]*TPL)
	} else {
		email.TPLs = conf.TPL
	}

	return email
}

func (e *Email) NewEmail() *email.Email {
	return email.NewEmail()
}

func (e *Email) NewEmailByTpl(tplName string, data interface{}) (*email.Email, error) {
	tpl, ok := e.conf.TPL[tplName]
	if !ok {
		return nil, fmt.Errorf("tpl name %s is not defined", tplName)
	}
	var err error
	mail := email.NewEmail()
	cachePrefix := e.cacheName + "." + tplName + "."
	subject, err := templateParse(cachePrefix+"subject", tpl.Subject, data)
	if err != nil {
		return nil, err
	}
	mail.Subject = string(subject)

	if tpl.Text != "" {
		mail.Text, err = templateParse(cachePrefix+"text", tpl.Text, data)
		if err != nil {
			return nil, err
		}
	}
	if tpl.HTML != "" {
		mail.HTML, err = templateParse(cachePrefix+"html", tpl.HTML, data)
		if err != nil {
			return nil, err
		}
	}

	mail.From = fmt.Sprintf("%s <%s>", e.conf.Name, e.conf.Account)
	if len(tpl.Bcc) > 0 {
		mail.Bcc = tpl.Bcc
	}

	if len(tpl.Cc) > 0 {
		mail.Cc = tpl.Cc
	}

	if len(tpl.To) > 0 {
		mail.To = tpl.To
	} else {
		mail.To = []string{}
	}

	return mail, nil
}

func (e *Email) Send(addr string, mail *email.Email) error {
	e.sendMu.Lock()
	defer e.sendMu.Unlock()

	mail.To = append(mail.To, addr)
	if e.TLS == nil {
		return mail.Send(e.mailHostAddr, e.SmtpAuth)
	}
	return mail.SendWithTLS(e.mailHostAddr, e.SmtpAuth, e.TLS)
}

func (e *Email) SendByTpl(addr string, tplName string, data interface{}) error {
	mail, err := e.NewEmailByTpl(tplName, data)
	if err != nil {
		return err
	}
	return e.Send(addr, mail)
}

var templateCache = map[string]*template.Template{}
var cacheMu sync.RWMutex

func templateParse(name, src string, data interface{}) ([]byte, error) {
	var err error
	cacheMu.RLock()
	parse, ok := templateCache[name]
	cacheMu.RUnlock()
	if !ok {
		parse, err = template.New(name).Parse(src)
		if err != nil {
			return nil, err
		}
		cacheMu.Lock()
		templateCache[name] = parse
		cacheMu.Unlock()
	}
	content := new(bytes.Buffer)
	err = parse.Execute(content, data)
	if err != nil {
		return nil, err
	}
	return content.Bytes(), nil
}
