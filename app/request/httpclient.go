package request

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	// HTTPClient http request client
	HTTPClient struct {
		http.Client
		transport *http.Transport // Client transport的引用
		UserAgent string
	}
)

// NewHTTPClient 返回 HTTPClient 指针
func NewHTTPClient() *HTTPClient {
	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}
	c := http.Client{
		Timeout:   30 * time.Second,
		Transport: t,
	}

	return &HTTPClient{Client: c, transport: t}
}

// SetProxy 设置代理
func (c *HTTPClient) SetProxy(addr string) (err error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	u := &url.URL{Host: net.JoinHostPort(host, port)}
	c.transport.Proxy = http.ProxyURL(u)
	return
}

// SetEnvProxy 设置为系统环境代理
func (c *HTTPClient) SetEnvProxy() {
	c.transport.Proxy = http.ProxyFromEnvironment
}

// SetUserAgent 设置浏览器标识
func (c *HTTPClient) SetUserAgent(ua string) {
	c.UserAgent = ua
}

// SetTimeout 设置请求超时时间
func (c *HTTPClient) SetTimeout(t time.Duration) {
	c.Client.Timeout = t
}

// Req 封装HTTP请求方法
func (c *HTTPClient) Req(method string, urlStr string, post interface{}, contentType string, header map[string]string) (resp *http.Response, err error) {
	var (
		req  *http.Request
		body io.Reader
	)

	if post != nil {
		switch data := post.(type) {
		case io.Reader:
			body = data
		case map[string]string:
			query := url.Values{}
			for k, v := range data {
				query.Set(k, v)
			}
			body = strings.NewReader(query.Encode())
		case map[string]interface{}:
			query := url.Values{}
			for k, v := range data {
				query.Set(k, fmt.Sprint(v))
			}
			body = strings.NewReader(query.Encode())
		case map[interface{}]interface{}:
			query := url.Values{}
			for k, v := range data {
				query.Set(fmt.Sprint(k), fmt.Sprint(v))
			}
			body = strings.NewReader(query.Encode())
		case string:
			body = strings.NewReader(data)
		case []byte:
			body = bytes.NewReader(data)
		default:
			return nil, fmt.Errorf("unknown post type:%s", data)
		}
	}

	req, err = http.NewRequest(method, urlStr, body)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", c.UserAgent)

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	return c.Client.Do(req)
}

// Fetch 获取页面主体
func (c *HTTPClient) Fetch(method string, urlStr string, post interface{}, contentType string, header map[string]string) (body []byte, err error) {
	resp, err := c.Req(method, urlStr, post, contentType, header)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
