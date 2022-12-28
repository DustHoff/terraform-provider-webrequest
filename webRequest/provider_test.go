package webRequest

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"testing"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"webrequest": providerserver.NewProtocol6WithError(NewProvider("test")()),
}

func testAccPreCheck(t *testing.T) {
	testAccProtoV6ProviderFactories["webrequest"]()
}
