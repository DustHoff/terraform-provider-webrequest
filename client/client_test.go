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
	result := testClient.NewRequest().SetMethod("FETCH").SetURL("https://httpbin.org/").Do()

	if result.StatusCode() != 405 {
		t.Fatal("StatusCode not 405(" + fmt.Sprintf("%v", result.statusCode) + ")")
	}
}

func TestClientSendGETInvalidURL(t *testing.T) {
	t.Log("Send Request to empty url")
	result := testClient.NewRequest().SetURL("").SetMethod("GET").Do()

	if result.StatusCode() != 499 {
		t.Fatal("StatusCode not 499(" + fmt.Sprintf("%v", result.statusCode) + ")")
	}

}

func TestClientSendGETEmptyBody(t *testing.T) {
	t.Log("Send GET request to https://httpbin.org/get")
	result := testClient.NewRequest().SetURL("https://httpbin.org/get").SetMethod("GET").Do()

	if result.StatusCode() != 200 {
		t.Fatal("Status Code not 200")
	}
}

func TestClientSendPostBody(t *testing.T) {
	t.Log("Send POST request to https://httpbin.org/post")
	result := testClient.NewRequest().SetURL("https://httpbin.org/post").SetMethod("POST").SetBody("body").Do()

	if result.StatusCode() != 200 {
		t.Fatal("Status Code not 200")
	}
	if len(result.Body()) == 0 {
		t.Fatal("received empty Response")
	}
	if result.BodyToJSON()["data"].(string) != "body" {
		t.Log(fmt.Sprintf("%v", result.Body()))
		t.Fatal("Response not parseable")
	}
}
