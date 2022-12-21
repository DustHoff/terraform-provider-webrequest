package webRequest

import (
	"context"
	"curl-terraform-provider/client"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWebRequest() *schema.Resource {
	return &schema.Resource{
		ReadContext: sendRequest,
		Schema: map[string]*schema.Schema{
			"result": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Response body of the requested resource",
			},
			"expires": &schema.Schema{
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Unix Timestamp",
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "time to live about the received response",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL of the target service. It includes schema, hostname, port and context path",
			},
			"body": &schema.Schema{
				Type:        schema.TypeString,
				Default:     nil,
				Optional:    true,
				Description: "Request Body for the request. please keep in mind to set the content-type header",
			},
			"method": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "GET",
				Optional:    true,
				Description: "The request method.",
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

func sendRequest(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if d.Get("expires").(int64) < time.Now().Unix() && d.Get("result") != nil {
		return diags
	}
	headers := d.Get("header").([]interface{})
	request := client.NewRequest().SetMethod(d.Get("method").(string)).SetURL(d.Get("url").(string))
	for _, entry := range headers {
		element := entry.(map[string]interface{})
		request.AddHeader(element["name"].(string), element["value"].(string))
	}
	res := request.Do()
	if res.StatusCode() != 200 {
		return diags
	}
	d.Set("result", res.Body())

	//set the expires timestamp
	if d.Get("ttl").(int) > 0 {
		d.Set("expires", time.Now().Unix()+d.Get("ttl").(int64))
	} else {
		d.Set("expires", time.Now().Unix())
	}
	// force that it always sets for the newest json object by changing the id of the object
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
