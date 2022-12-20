package client

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	statusCode int
	body       string
}

func NewResponse(res http.Response) Response {
	body, _ := io.ReadAll(res.Body)
	return Response{
		statusCode: res.StatusCode,
		body:       string(body),
	}
}

func NewCustomResponse(statusCode int, body string) Response {
	return Response{
		statusCode: statusCode,
		body:       body,
	}
}

func (r Response) StatusCode() int {
	return r.statusCode
}

func (r Response) Body() string {
	return r.body
}

func (r Response) BodyToJSON() map[string]interface{} {
	var parsed map[string]interface{}
	json.Unmarshal([]byte(r.body), &parsed)
	return parsed
}
