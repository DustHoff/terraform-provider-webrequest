package webRequest

import (
	"context"
	"curl-terraform-provider/client"
	"curl-terraform-provider/helper/modifier"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &RestDataCall{}

type RestDataCall struct {
	client client.Client
}

type RestDataCallModel struct {
	ID               types.String `tfsdk:"id"`
	Key              types.String `tfsdk:"key"`
	Result           types.String `tfsdk:"result"`
	StatusCode       types.Int64  `tfsdk:"statuscode"`
	URL              types.String `tfsdk:"url"`
	Body             types.String `tfsdk:"body"`
	IgnoreStatusCode types.Bool   `tfsdk:"ignorestatuscode"`
	Header           types.Map    `tfsdk:"header"`
	Create           types.Object `tfsdk:"create"`
	Read             types.Object `tfsdk:"read"`
	Update           types.Object `tfsdk:"update"`
	Delete           types.Object `tfsdk:"delete"`
}

type CustomAPICall struct {
	Method types.String `tfsdk:"method"`
	URL    types.String `tfsdk:"url"`
}

func NewRestDataCall() resource.Resource {
	return &RestDataCall{}
}

func (r *RestDataCall) sendRequest(ctx context.Context, data RestDataCallModel, custom CustomAPICall) client.Response {
	request := r.client.NewRequest().SetMethod(custom.Method.ValueString()).SetURL(custom.URL.ValueString()).SetBody(data.Body.ValueString())
	tflog.Debug(ctx, "Header Count "+fmt.Sprint(len(data.Header.Elements())))
	for key, value := range data.Header.Elements() {
		tflog.Debug(ctx, "Adding Header "+key+"="+value.(types.String).ValueString())
		request.AddHeader(key, value.(types.String).ValueString())
	}
	return request.Do()
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
			"statuscode": schema.Int64Attribute{
				Computed: true,
			},
			"ignorestatuscode": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					modifier.BooleanDefaultValueModifier(types.BoolValue(false)),
				},
			},
			"url": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"body": schema.StringAttribute{
				Optional: true,
			},
			"key": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					modifier.StringDefaultValue(types.StringValue("id")),
				},
			},
			"header": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					modifier.EmptyMapDefaultValue(),
				},
			},
			"create": schema.ObjectAttribute{
				Optional: true,
				Computed: true,
				AttributeTypes: map[string]attr.Type{
					"method": types.StringType,
					"url":    types.StringType,
				},
				PlanModifiers: []planmodifier.Object{
					modifier.CustomCallDefaultValueModifier("POST"),
				},
			},
			"read": schema.ObjectAttribute{
				Optional: true,
				Computed: true,
				AttributeTypes: map[string]attr.Type{
					"method": types.StringType,
					"url":    types.StringType,
				},
				PlanModifiers: []planmodifier.Object{
					modifier.CustomCallDefaultValueModifier("GET"),
				},
			},
			"update": schema.ObjectAttribute{
				Optional: true,
				Computed: true,
				AttributeTypes: map[string]attr.Type{
					"method": types.StringType,
					"url":    types.StringType,
				}, PlanModifiers: []planmodifier.Object{
					modifier.CustomCallDefaultValueModifier("PUT"),
				},
			},
			"delete": schema.ObjectAttribute{
				Optional: true,
				Computed: true,
				AttributeTypes: map[string]attr.Type{
					"method": types.StringType,
					"url":    types.StringType,
				},
				PlanModifiers: []planmodifier.Object{
					modifier.CustomCallDefaultValueModifier("DELETE"),
				},
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
	var custom CustomAPICall

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	customDiag := data.Create.As(ctx, &custom, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: false,
	})

	if customDiag != nil {
		resp.Diagnostics.Append(customDiag...)
	}
	if custom.URL.IsNull() {
		custom.URL = data.URL
	}
	response := r.sendRequest(ctx, data, custom)
	data.StatusCode = types.Int64Value(int64(response.StatusCode()))
	if (data.IgnoreStatusCode.ValueBool()) || (response.StatusCode() == 200) {
		data.Result = types.StringValue(response.Body())
		data.ID = types.StringValue(fmt.Sprint(response.BodyToJSON()[data.Key.ValueString()]))
	} else {
		resp.Diagnostics.AddError("Failed to create data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RestDataCall) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RestDataCallModel
	var custom CustomAPICall

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	customDiag := data.Read.As(ctx, &custom, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: false,
	})

	if customDiag != nil {
		resp.Diagnostics.Append(customDiag...)
	}
	if custom.URL.IsNull() {
		custom.URL = types.StringValue(data.URL.ValueString() + "/" + data.ID.ValueString())
	}
	response := r.sendRequest(ctx, data, custom)
	if (data.IgnoreStatusCode.ValueBool()) || (response.StatusCode() == 200) {
		data.Result = types.StringValue(response.Body())
	} else {
		resp.Diagnostics.AddError("Failed to fetch data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RestDataCall) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RestDataCallModel
	var custom CustomAPICall

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	customDiag := data.Update.As(ctx, &custom, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: false,
	})

	if customDiag != nil {
		resp.Diagnostics.Append(customDiag...)
	}
	if custom.URL.IsNull() {
		custom.URL = types.StringValue(data.URL.ValueString() + "/" + data.ID.ValueString())
	}
	response := r.sendRequest(ctx, data, custom)
	data.StatusCode = types.Int64Value(int64(response.StatusCode()))
	if (data.IgnoreStatusCode.ValueBool()) || (response.StatusCode() == 200) {
		data.Result = types.StringValue(response.Body())
	} else {
		resp.Diagnostics.AddError("Failed to update data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RestDataCall) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RestDataCallModel
	var custom CustomAPICall

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	customDiag := data.Delete.As(ctx, &custom, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: false,
	})

	if customDiag != nil {
		resp.Diagnostics.Append(customDiag...)
	}
	if custom.Method.IsNull() {
		custom.Method = types.StringValue("DELETE")
	}
	if custom.URL.IsNull() {
		custom.URL = types.StringValue(data.URL.ValueString() + "/" + data.ID.ValueString())
	}
	response := r.sendRequest(ctx, data, custom)
	if (data.IgnoreStatusCode.ValueBool()) || (response.StatusCode() == 200) {
		data.Result = types.StringValue(response.Body())
	} else {
		resp.Diagnostics.AddError("Failed to fetch data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
}
