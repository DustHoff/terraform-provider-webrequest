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
	url = "https://eoscet74ykdzldt.m.pipedream.net"
	ignorestatuscode = true
	body = jsonencode({"username":"test","email":"test@example.com"})
	header = {
		Content-Type = "application/json"
		Accept = "application/json"
	}

	create = {
		method = "POST"
		url = "https://eoscet74ykdzldt.m.pipedream.net/create"
	}

	read = {
		method = "POST"
		url = "https://eoscet74ykdzldt.m.pipedream.net/get"
	}

	update = {
		method = "POST"
		url = "https://eoscet74ykdzldt.m.pipedream.net/update"
	}

	delete = {
		method = "POST"
		url = "https://eoscet74ykdzldt.m.pipedream.net/delete"
	}
}
`
