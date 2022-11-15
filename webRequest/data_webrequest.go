package webRequest

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWebRequest() *schema.Resource {
	return &schema.Resource{
		ReadContext: sendRequest,
		Schema: map[string]*schema.Schema{
			"result": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"URL": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"body": &schema.Schema{
				Type:     schema.TypeString,
				Default:  nil,
				Optional: true,
			},
			"method": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "GET",
				Optional: true,
			},
			"header": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func sendRequest(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 60 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest(d.Get("method").(string), d.Get("URL").(string), strings.NewReader(d.Get("body").(string)))
	headers := d.Get("header").([]interface{})
	for _, entry := range headers {
		element := entry.(map[string]interface{})
		req.Header.Set(element["name"].(string), element["value"].(string))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	responseString := string(responseData)

	//set the actual value we're going to return into, associated with the 'response' key name.
	if err := d.Set("result", responseString); err != nil {
		return diag.FromErr(err)
	}

	// force that it always sets for the newest json object by changing the id of the object
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
