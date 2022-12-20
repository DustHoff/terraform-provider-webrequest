package webRequest

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

var providerConfiguration map[string]*schema.Provider
var testProvider *schema.Provider

func init() {
	testProvider = Provider()
	providerConfiguration = map[string]*schema.Provider{
		"hashicups": testProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := testProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
