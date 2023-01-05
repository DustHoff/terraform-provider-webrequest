package webRequest

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestComplexWebRequestDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: complexWebrequestSend,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.webrequest_send.test", "id"),
					resource.TestCheckResourceAttrSet("data.webrequest_send.test", "result"),
					resource.TestCheckResourceAttrSet("data.webrequest_send.test", "expires"),
				),
			},
		},
	})
}

const complexWebrequestSend = `
data "webrequest_send" "test" {
  	url = "https://httpbin.org/post"
  	method = "POST"
  	body = jsonencode({"username":"test","email":"test@example.com"})
  	header = {
			Content-Type = "application/json"
			Accept = "application/json"
  	}
}
`
