package webRequest

import (
	"context"
	"curl-terraform-provider/client"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &RestDataCall{}

type RestDataCall struct {
	client client.Client
}

type RestDataCallModel struct {
	ID     types.String `tfsdk:"id"`
	Key    types.String `tfsdk:"key"`
	Result types.String `tfsdk:"result"`
	URL    types.String `tfsdk:"url"`
	Body   types.String `tfsdk:"body"`
	Header types.Map    `tfsdk:"header"`
}

func NewRestDataCall() resource.Resource {
	return &RestDataCall{}
}

func (r *RestDataCall) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_restcall"
}

func (r *RestDataCall) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"result": schema.StringAttribute{
				Computed: true,
			},
			"url": schema.StringAttribute{
				Required: true,
			},
			"body": schema.StringAttribute{
				Optional: true,
			},
			"key": schema.StringAttribute{
				Optional: true,
			},
			"header": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *RestDataCall) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *RestDataCall) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RestDataCallModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	request := r.client.NewRequest().SetMethod("POST").SetURL(data.URL.ValueString()).SetBody(data.Body.ValueString())
	if !data.Header.IsUnknown() {
		for key, value := range data.Header.Elements() {
			request.AddHeader(key, value.String())
		}
	}
	response := request.Do()
	if (response.StatusCode() >= 200) && (response.StatusCode() < 299) {
		data.Result = types.StringValue(response.Body())
		data.ID = types.StringValue(fmt.Sprint(response.BodyToJSON()[data.Key.ValueString()]))
	} else {
		resp.Diagnostics.AddError("Failed to create data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RestDataCall) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RestDataCallModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	request := r.client.NewRequest().SetMethod("GET").SetURL(data.URL.ValueString() + "/" + data.ID.ValueString())
	for key, value := range data.Header.Elements() {
		request.AddHeader(key, value.String())
	}
	response := request.Do()
	if (response.StatusCode() >= 200) && (response.StatusCode() < 299) {
		data.Result = types.StringValue(response.Body())
	} else {
		resp.Diagnostics.AddError("Failed to fetch data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RestDataCall) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RestDataCallModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	request := r.client.NewRequest().SetMethod("PUT").SetURL(data.URL.ValueString() + "/" + data.ID.ValueString())
	for key, value := range data.Header.Elements() {
		request.AddHeader(key, value.String())
	}
	response := request.Do()
	if (response.StatusCode() >= 200) && (response.StatusCode() < 299) {
		data.Result = types.StringValue(response.Body())
	} else {
		resp.Diagnostics.AddError("Failed to update data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RestDataCall) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RestDataCallModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	request := r.client.NewRequest().SetMethod("DELETE").SetURL(data.URL.ValueString() + "/" + data.ID.ValueString())
	for key, value := range data.Header.Elements() {
		request.AddHeader(key, value.String())
	}
	response := request.Do()
	if (response.StatusCode() >= 200) && (response.StatusCode() < 299) {
		data.Result = types.StringValue(response.Body())
	} else {
		resp.Diagnostics.AddError("Failed to fetch data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
}
