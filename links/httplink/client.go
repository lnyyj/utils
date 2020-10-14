package httplink

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPC http客户端对象
type HTTPC struct {
	Cli        http.Client
	streamFlag bool
}

// NewHTTPC 创建一个http客户端对象
func NewHTTPC(maxconns ...int) *HTTPC {
	h := &HTTPC{}
	h.Cli = http.Client{}
	h.SetStreamFlag(false)
	if len(maxconns) == 1 && maxconns[0] > 0 {
		mineTransport := &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			MaxIdleConnsPerHost: maxconns[0], //连接池个数
		}

		h.Cli.Transport = mineTransport
		h.Cli.Timeout = 50 * time.Second
		h.SetStreamFlag(true)
	}

	return h
}

// NewHTTPSC 创建一个https客户端对象
func NewHTTPSC(maxconns int, crtpath ...string) *HTTPC {
	h := &HTTPC{}
	h.Cli = http.Client{}
	pool := x509.NewCertPool()
	h.SetStreamFlag(true)

	// 可添加多个证书
	for _, v := range crtpath {
		caCrt, err := ioutil.ReadFile(v)
		if err != nil {
			panic("load crt: " + err.Error())
		}

		pool.AppendCertsFromPEM(caCrt)
	}

	h.Cli.Timeout = 50 * time.Second
	h.Cli.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		MaxIdleConnsPerHost: maxconns, //连接池个数
		TLSClientConfig:     &tls.Config{RootCAs: pool, InsecureSkipVerify: true},
	}

	return h
}

// SetStreamFlag 客户端标记
func (h *HTTPC) SetStreamFlag(v bool) *HTTPC {
	h.streamFlag = v
	return h
}

// NewHeader 创建一个头部
func (h *HTTPC) NewHeader(kv ...string) (header http.Header) {
	l := len(kv)
	if l%2 != 0 {
		panic("new header num fail")
	}

	header = make(http.Header)
	for k := 0; k < l; k += 2 {
		header.Add(kv[k], kv[k+1])
	}

	return header
}

// SendUpstream 发送数据
// http 客户端发送数据
// packet 要发送报文
// method 选择发送方式 post get
// url 地址
func (h *HTTPC) SendUpstream(header http.Header, packet, method, url string) ([]byte, error) {
	// if !h.streamFlag {
	// 	return nil, errors.New("init http client error")
	// }
	request, err := http.NewRequest(method, url, strings.NewReader(packet))
	if err != nil {
		return nil, errors.New("NewRequest: " + err.Error())
	}
	request.Header = header

	response, err := h.Cli.Do(request)
	if err != nil {
		return nil, errors.New("Do: " + err.Error())
	}
	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("ReadAll: " + err.Error())
	}

	return buf, nil
}

// HTTPSPostForm post表单
func (h *HTTPC) HTTPSPostForm(url string, data url.Values) ([]byte, error) {
	var resp *http.Response
	var err error

	if h.streamFlag {
		resp, err = h.Cli.PostForm(url, data)
	} else {
		resp, err = http.PostForm(url, data)
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// HTTPGet get表单
func (h *HTTPC) HTTPGet(url string) ([]byte, error) {
	var resp *http.Response
	var err error

	if h.streamFlag {
		resp, err = h.Cli.Get(url)
	} else {
		resp, err = http.Get(url)
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetCookie 返回Cookie
func (h *HTTPC) GetCookie(header http.Header, packet, method, url string) (string, error) {
	if !h.streamFlag {
		return "", errors.New("init http client error")
	}

	request, err := http.NewRequest(method, url, strings.NewReader(packet))
	if err != nil {
		return "", errors.New("NewRequest: " + err.Error())
	}
	request.Header = header

	response, err := h.Cli.Do(request)
	if err != nil {
		return "", errors.New("Do: " + err.Error())
	}
	defer response.Body.Close()
	Cookie := ""
	if response.Header != nil {
		if _, ok := response.Header["Set-Cookie"]; ok {
			Cookie = response.Header["Set-Cookie"][0]
		}
	}
	return Cookie, nil
}
