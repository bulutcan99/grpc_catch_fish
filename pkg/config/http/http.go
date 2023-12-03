package config_http

import (
	"errors"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

type HttpClient struct {
	client *fasthttp.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{client: &fasthttp.Client{
		ReadTimeout:         5 * time.Second,
		WriteTimeout:        5 * time.Second,
		MaxIdleConnDuration: 5 * time.Second,
		MaxConnWaitTimeout:  30 * time.Second,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}}
}

func (h *HttpClient) IsHostAlive(url string) bool {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	err := h.client.Do(req, resp)

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if err != nil {
		return false
	}
	return resp.StatusCode() == 200
}

func (h *HttpClient) Get(url string) ([]byte, error) {
	statusCode, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return nil, err
	}

	if statusCode != fasthttp.StatusOK {
		return nil, errors.New("failed to get data, status code: " + strconv.Itoa(statusCode))
	}

	return body, nil
}
