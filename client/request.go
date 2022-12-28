package client

import (
	"net/http"
	"strings"
)

type Request struct {
	client Client
	url    string
	method string
	body   string
	header map[string]string
}

func NewRequestBuilder(client Client) *Request {
	return &Request{
		client: client,
	}
}

func (r *Request) AddHeader(key string, value string) *Request {
	r.header[key] = value
	return r
}

func (r *Request) SetURL(url string) *Request {
	r.url = url
	return r
}
func (r *Request) SetMethod(method string) *Request {
	r.method = method
	return r
}

func (r *Request) SetBody(body string) *Request {
	r.body = body
	return r
}
func (r *Request) Do() Response {
	request, err := http.NewRequest(r.method, r.url, strings.NewReader(r.body))
	for key, value := range r.header {
		request.Header.Set(key, value)
	}
	if err != nil {
		return NewCustomResponse(444, err.Error())
	}
	return r.client.do(request)
}
