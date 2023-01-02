package webRequest

import (
	"context"
	"curl-terraform-provider/client"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
		MarkdownDescription: "send any http request to a endpoint",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "data source identifier",
				Computed:            true,
			},
			"result": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "received response data as string",
			},
			"expires": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "timestamp when the received response invalidates. value is in milliseconds since 1970",
			},
			"ttl": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "time to live of the received response. value is in seconds",
			},
			"url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "request url",
			},
			"body": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "request body",
			},
			"method": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "request method",
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("GET", "POST", "PUT", "DELETE", "OPTION", "HEAD"),
				},
			},
			"header": schema.MapAttribute{
				Optional:            true,
				MarkdownDescription: "map of request header",
				ElementType:         types.StringType,
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

	tflog.Debug(ctx, "Header Count "+fmt.Sprint(len(data.Header.Elements())))
	for key, value := range data.Header.Elements() {
		tflog.Debug(ctx, "Adding Header "+key+"="+value.(types.String).ValueString())
		request.AddHeader(key, value.(types.String).ValueString())
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
