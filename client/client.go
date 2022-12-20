package client

import (
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(timeoutSecs int) Client {
	client := Client{
		httpClient: &http.Client{Timeout: time.Duration(timeoutSecs) * time.Second},
	}
	return client
}

func (c Client) do(req *http.Request) Response {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NewCustomResponse(499, err.Error())
	}

	return NewResponse(*resp)
}

func (c Client) Send(method string, url string, body string) Response {
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	return c.do(req)
}
