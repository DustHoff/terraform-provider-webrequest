package webRequest

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestComplexResourceRestDataCall(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: complexRestCall,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("webrequest_restcall.call", "id", "test"),
					resource.TestCheckResourceAttr("webrequest_restcall.call", "result", "{\"email\":\"test@example.com\",\"username\":\"test\"}"),
					resource.TestCheckResourceAttrSet("webrequest_restcall.call", "result"),
					resource.TestCheckResourceAttrSet("webrequest_restcall.call", "id"),
				),
			},
		},
	})
}

const complexRestCall = `
resource "webrequest_restcall" "call" {
	body = jsonencode({"username":"test","email":"test@example.com"})
	url = "https://httpbin.org/post"
	ignorestatuscode = true
	filter = "//data"
	key = "//json/username"
	header = {
		Content-Type = "application/json"
		Accept = "application/json"
	}

	create = {
		method = "POST"
	}

	read = {
		method = "PATCH"
		url = "https://httpbin.org/patch"
		body = jsonencode({"id":"{ID}","username":"test","email":"test@example.com"})
	}

	update = {
		method = "PUT"
		url = "https://httpbin.org/put"
		filter = "//data"
		body = jsonencode({"username":"test","email":"test@example.com"})
	}

	delete = {
		method = "DELETE"
		url = "https://httpbin.org/delete"
		body = jsonencode({"username":"test","email":"test@example.com"})
	}

  	lifecycle {
    	postcondition {
      	condition     = self.statuscode == 200
      	error_message = "Received Statuscode should be http/200"
    	}
  	}
}
`
