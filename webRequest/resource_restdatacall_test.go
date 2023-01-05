package webRequest

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestSimpleResourceRestDataCall(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: simpleRestCall,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("webrequest_restcall.call", "id", "1"),
					resource.TestCheckResourceAttrSet("webrequest_restcall.call", "result"),
					resource.TestCheckResourceAttrSet("webrequest_restcall.call", "id"),
				),
			},
		},
	})
}

func TestComplexResourceRestDataCall(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: complexRestCall,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("webrequest_restcall.call", "id", "1"),
					resource.TestCheckResourceAttrSet("webrequest_restcall.call", "result"),
					resource.TestCheckResourceAttrSet("webrequest_restcall.call", "id"),
				),
			},
		},
	})
}

const simpleRestCall = `
resource "webrequest_restcall" "call" {
	url = "https://eoscet74ykdzldt.m.pipedream.net"
	body = jsonencode({"username":"test","email":"test@example.com"})
}
`

const complexRestCall = `
resource "webrequest_restcall" "call" {
	ignorestatuscode = true
	header = {
		Content-Type = "application/json"
		Accept = "application/json"
	}

	create = {
		method = "POST"
		url = "https://eoscet74ykdzldt.m.pipedream.net/create"
		body = jsonencode({"username":"test","email":"test@example.com"})
	}

	read = {
		method = "POST"
		url = "https://eoscet74ykdzldt.m.pipedream.net/get/{ID}"
		body = jsonencode({"id":"{ID}","username":"test","email":"test@example.com"})
	}

	update = {
		method = "POST"
		url = "https://eoscet74ykdzldt.m.pipedream.net/update"
		body = jsonencode({"username":"test","email":"test@example.com"})
	}

	delete = {
		method = "POST"
		url = "https://eoscet74ykdzldt.m.pipedream.net/delete"
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
