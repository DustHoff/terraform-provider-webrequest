package client

import (
	"fmt"
	"testing"
)

var testClient Client

func init() {
	testClient = NewClient(30)
}
func TestClientSendInvalidMethod(t *testing.T) {
	t.Log("Send Request with Method Fetch to https://httpbin.org/")
	result := testClient.Send("FETCH", "https://httpbin.org/", "")

	if result.statusCode != 405 {
		t.Fatal("StatusCode not 405(" + fmt.Sprintf("%v", result.statusCode) + ")")
	}
}

func TestClientSendGETInvalidURL(t *testing.T) {
	t.Log("Send Request to empty url")
	result := testClient.Send("GET", "", "")

	if result.statusCode != 499 {
		t.Fatal("StatusCode not 499(" + fmt.Sprintf("%v", result.statusCode) + ")")
	}

}

func TestClientSendGETEmptyBody(t *testing.T) {
	t.Log("Send GET request to https://httpbin.org/get")
	result := testClient.Send("GET", "https://httpbin.org/get", "")

	if result.statusCode != 200 {
		t.Fatal("Status Code not 200")
	}
}

func TestClientSendPostBody(t *testing.T) {
	t.Log("Send POST request to https://httpbin.org/post")
	result := testClient.Send("POST", "https://httpbin.org/post", "body")

	if result.statusCode != 200 {
		t.Fatal("Status Code not 200")
	}
	if len(result.Body()) == 0 {
		t.Fatal("received empty Response")
	}
	if result.BodyToJSON()["data"].(string) != "body" {
		t.Fatal("Response not parseable")
	}
}
