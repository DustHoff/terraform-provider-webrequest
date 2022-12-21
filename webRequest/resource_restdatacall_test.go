package webRequest

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestResourceRestDataCall(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"webrequest": testProvider,
		},
		Steps: []resource.TestStep{
			{
				Config: generateSimpleTerraformConfigRestDataCall(),
				Check: resource.ComposeTestCheckFunc(
					checkTerraformState("webrequest_restcall.call", t),
				),
			},
		},
	})
}

func generateSimpleTerraformConfigRestDataCall() string {
	return fmt.Sprintf(`
	resource "webrequest_restcall" "call" {
		url = "https://eoscet74ykdzldt.m.pipedream.net"
		body = jsonencode({"username":"test","email":"test@example.com"})
	}
	`)
}

func checkTerraformState(call string, t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[call]

		if !ok {
			return fmt.Errorf("Not found: %s", call)
		}
		if rs.Primary.ID != "1" {
			return fmt.Errorf("id doesnt match")
		}

		return nil
	}
}
