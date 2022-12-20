package client

import (
	"net/http"
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

func (c Client) NewRequest() *Request {
	return NewRequestBuilder(c)
}
