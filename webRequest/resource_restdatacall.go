package webRequest

import (
	"context"
	"curl-terraform-provider/client"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRestDataCall() *schema.Resource {
	return &schema.Resource{
		CreateContext: createData,
		ReadContext:   fetchData,
		UpdateContext: updateData,
		DeleteContext: deleteData,
		Schema: map[string]*schema.Schema{
			"result": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Response body of the requested resource",
			},
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "id",
				Description: "Primary key of the response object",
			},
			"objectid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Primary key of the response object",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL of the target service. It includes schema, hostname, port and context path",
			},
			"body": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Request Body for the request. please keep in mind to set the content-type header",
			},
			"header": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of all Request Header",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Header name, like Content-Type",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Header value",
						},
					},
				},
			},
		},
	}
}

func createData(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*client.Client)
	request := client.NewRequest().SetMethod("POST").SetURL(d.Get("url").(string)).SetBody(d.Get("body").(string))
	headers := d.Get("header").([]interface{})
	for _, entry := range headers {
		element := entry.(map[string]interface{})
		request.AddHeader(element["name"].(string), element["value"].(string))
	}
	response := request.Do()
	if (response.StatusCode() >= 200) && (response.StatusCode() < 299) {
		d.Set("result", response.Body())
		d.Set("objectid", fmt.Sprint(response.BodyToJSON()[d.Get("key").(string)]))
		d.SetId(fmt.Sprint(response.BodyToJSON()[d.Get("key").(string)]))
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create data",
			Detail:   "received response code " + fmt.Sprint(response.StatusCode()),
		})
	}
	return diags
}

func fetchData(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*client.Client)
	request := client.NewRequest().SetMethod("GET").SetURL(d.Get("url").(string) + "/" + d.Get("objectid").(string))
	headers := d.Get("header").([]interface{})
	for _, entry := range headers {
		element := entry.(map[string]interface{})
		request.AddHeader(element["name"].(string), element["value"].(string))
	}
	response := request.Do()
	if (response.StatusCode() >= 200) && (response.StatusCode() < 299) {
		d.Set("result", response.Body())
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to fetch data",
			Detail:   "received response code " + fmt.Sprint(response.StatusCode()),
		})
	}
	return diags
}

func updateData(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*client.Client)
	request := client.NewRequest().SetMethod("PUT").SetURL(d.Get("url").(string) + "/" + d.Get("objectid").(string))
	headers := d.Get("header").([]interface{})
	for _, entry := range headers {
		element := entry.(map[string]interface{})
		request.AddHeader(element["name"].(string), element["value"].(string))
	}
	response := request.Do()
	if (response.StatusCode() >= 200) && (response.StatusCode() < 299) {
		d.Set("result", response.Body())
		d.SetId(response.BodyToJSON()[d.Get("key").(string)].(string))
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to update data",
			Detail:   "received response code " + fmt.Sprint(response.StatusCode()),
		})
	}
	return diags
}

func deleteData(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*client.Client)
	request := client.NewRequest().SetMethod("DELETE").SetURL(d.Get("url").(string) + "/" + d.Get("objectid").(string))
	headers := d.Get("header").([]interface{})
	for _, entry := range headers {
		element := entry.(map[string]interface{})
		request.AddHeader(element["name"].(string), element["value"].(string))
	}
	response := request.Do()

	if (response.StatusCode() >= 200) && (response.StatusCode() < 299) {
		d.SetId("")
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete data",
			Detail:   "received response code " + fmt.Sprint(response.StatusCode()),
		})
	}
	return diags
}
