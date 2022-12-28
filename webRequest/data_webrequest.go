package webRequest

import (
	"context"
	"curl-terraform-provider/client"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"time"
)

var _ datasource.DataSource = &WebRequestDataSource{}

func NewWebRequestDataSource() datasource.DataSource {
	return &WebRequestDataSource{}
}

type WebRequestDataSource struct {
	client client.Client
}
type WebRequestDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Result  types.String `tfsdk:"result"`
	Expires types.Int64  `tfsdk:"expires"`
	TTL     types.Int64  `tfsdk:"ttl"`
	URL     types.String `tfsdk:"url"`
	Body    types.String `tfsdk:"body"`
	Method  types.String `tfsdk:"method"`
	Header  types.Map    `tfsdk:"header"`
}

func (d *WebRequestDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_send"
}

func (d *WebRequestDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
			"result": schema.StringAttribute{
				Computed: true,
			},
			"expires": schema.Int64Attribute{
				Computed: true,
			},
			"ttl": schema.Int64Attribute{
				Optional: true,
			},
			"url": schema.StringAttribute{
				Required: true,
			},
			"body": schema.StringAttribute{
				Optional: true,
			},
			"method": schema.StringAttribute{
				Optional: true,
			},
			"header": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *WebRequestDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *WebRequestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WebRequestDataSourceModel
	tflog.Info(ctx, "read DataModel")
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	tflog.Info(ctx, "check for errors")
	if resp.Diagnostics.HasError() {
		return
	}

	if (data.Expires.ValueInt64() < time.Now().Unix()) && (!data.Result.IsNull()) {
		return
	}
	if data.URL.IsNull() || data.URL.IsUnknown() {
		resp.Diagnostics.AddError("URL not set", "No URL has been defined")
		return
	}
	request := d.client.NewRequest().SetMethod(data.Method.ValueString()).SetURL(data.URL.ValueString()).SetBody(data.Body.ValueString())

	if !data.Header.IsUnknown() {
		for key, value := range data.Header.Elements() {
			request.AddHeader(key, value.String())
		}
	}
	res := request.Do()
	if res.StatusCode() != 200 {
		resp.Diagnostics.AddError("Unhealthy Response Code "+fmt.Sprint(res.StatusCode()), res.Body())
	}
	data.Result = types.StringValue(res.Body())

	//set the expires timestamp
	if data.TTL.ValueInt64() > 0 {
		data.Expires = types.Int64Value(time.Now().Unix() + data.TTL.ValueInt64())
	} else {
		data.Expires = types.Int64Value(time.Now().Unix())
	}
	// force that it always sets for the newest json object by changing the id of the object
	data.ID = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
