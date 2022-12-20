package webRequest

import (
	"context"
	client2 "curl-terraform-provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"webrequest_restcall": resourceRestDataCall(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"webrequest_send": dataWebRequest(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := client2.NewClient(data.Get("timeout").(int))
	return client, diags
}
