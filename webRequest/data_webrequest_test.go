package webRequest

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestSimpleWebRequestDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: simpleWebrequestSend,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.webrequest_send.test", "id"),
					resource.TestCheckResourceAttrSet("data.webrequest_send.test", "result"),
					resource.TestCheckResourceAttrSet("data.webrequest_send.test", "expires"),
				),
			},
		},
	})
}

const simpleWebrequestSend = `
data "webrequest_send" "test" {
  url = "https://eoscet74ykdzldt.m.pipedream.net"
  method = "POST"
  body = jsonencode({"username":"test","email":"test@example.com"})
}
`

const complexWebrequestSend = `
data "webrequest_send" "test" {
  url = "https://eoscet74ykdzldt.m.pipedream.net"
  method = "POST"
  body = jsonencode({"username":"test","email":"test@example.com"})
  header {
    test = "test"
  }
}
`
