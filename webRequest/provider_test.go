package webRequest

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

var testProvider *schema.Provider

func init() {
	testProvider = Provider()
}

func TestProvider(t *testing.T) {
	if err := testProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
