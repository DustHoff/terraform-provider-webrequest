package webRequest

import (
	"context"
	"curl-terraform-provider/client"
	"curl-terraform-provider/helper/modifier"
	types2 "curl-terraform-provider/helper/types"
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"regexp"
	"strings"
)

var _ resource.ResourceWithConfigValidators = &RestDataCall{}

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
	Filter           types.String `tfsdk:"filter"`
	IgnoreStatusCode types.Bool   `tfsdk:"ignorestatuscode"`
	Header           types.Map    `tfsdk:"header"`
	Create           types.Object `tfsdk:"create"`
	Read             types.Object `tfsdk:"read"`
	Update           types.Object `tfsdk:"update"`
	Delete           types.Object `tfsdk:"delete"`
}

func NewRestDataCall() resource.Resource {
	return &RestDataCall{}
}

func (r *RestDataCall) sendRequest(ctx context.Context, data RestDataCallModel, custom types2.AlternativeRequestParameter) client.Response {
	method := custom.GetMethod()
	url := custom.GetURL()
	body := custom.GetBody()
	tflog.Info(ctx, "Configured URL "+data.URL.ValueString())
	tflog.Info(ctx, "Configured Custom URL "+custom.GetURL().ValueString())
	if url.IsNull() {
		url = types.StringValue(data.URL.ValueString())
	}

	tflog.Info(ctx, "Configured Body "+data.Body.ValueString())
	tflog.Info(ctx, "Configured Custom Body "+body.ValueString())
	if body.IsNull() {
		body = data.Body
	}

	regex := regexp.MustCompile("{ID}")

	tflog.Info(ctx, ">> "+method.ValueString()+" "+regex.ReplaceAllString(url.ValueString(), data.ID.ValueString()))
	tflog.Info(ctx, ">> "+regex.ReplaceAllString(body.ValueString(), data.ID.ValueString()))
	request := r.client.NewRequest().SetMethod(method.ValueString()).SetURL(regex.ReplaceAllString(url.ValueString(), data.ID.ValueString())).SetBody(regex.ReplaceAllString(body.ValueString(), data.ID.ValueString()))

	for key, value := range data.Header.Elements() {
		tflog.Info(ctx, ">> "+"Adding Header "+key+"="+value.(types.String).ValueString())
		request.AddHeader(key, value.(types.String).ValueString())
	}
	return request.Do()
}

func (r *RestDataCall) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_restcall"
}

func (r *RestDataCall) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *RestDataCall) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource interact with any rest like endpoint. All CRUD types are handled during lifetime of " +
			"this resource. fresh resource generate a create request, refreshing resource state result in a read action, " +
			"updating a attribute perform after apply a update action (partial update isn't supported) and finally deleting the resource performs a delete request " +
			"the primary key ob the resulting object is append to the url " +
			"You can use the placeholer {ID} in URL and Body attribute to append or add the resource id(primary key value) to your request",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "received object identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"result": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "received JSON object as string representation",
			},
			"statuscode": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "received http statuscode",
			},
			"ignorestatuscode": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "ignores the statuscode on response validation",
				PlanModifiers: []planmodifier.Bool{
					modifier.BooleanDefaultValueModifier(types.BoolValue(false)),
				},
			},
			"url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "request url",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"body": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "request body",
			},
			"filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "JSON path expression to filter selective the value of attribute result. The expression based on XPath",
			},
			"key": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "primary key of the received object, to generate/manipulate the request url, use a JSON path expression to get the right value. The expression based on XPath ",
				PlanModifiers: []planmodifier.String{
					modifier.StringDefaultValue(types.StringValue("id")),
				},
			},
			"header": schema.MapAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "map of request header",
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					modifier.EmptyMapDefaultValue(),
				},
			},
			"create": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "manipulate the behavior for object creation",
				Attributes: map[string]schema.Attribute{
					"method": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "HTTP Method to use for the request",
						PlanModifiers: []planmodifier.String{
							modifier.StringDefaultValue(types.StringValue("POST")),
						},
					},
				},
			},
			"read": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "manipulate the behavior for reading the object",
				Attributes: map[string]schema.Attribute{
					"method": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "HTTP Method to use for the request",
						PlanModifiers: []planmodifier.String{
							modifier.StringDefaultValue(types.StringValue("GET")),
						},
					},
					"url": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Request URL to read the JSON Object.",
					},
					"body": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Custom request body for the request",
					},
					"filter": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Custom result filter for the request",
					},
					"keepid": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Flag to keep the initial resource id, even when the response contains a new one",
					},
				},
			},
			"update": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "manipulate the behavior for updating the object",
				Attributes: map[string]schema.Attribute{
					"method": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "HTTP Method to use for the request",
						PlanModifiers: []planmodifier.String{
							modifier.StringDefaultValue(types.StringValue("PUT")),
						},
					},
					"url": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Request URL to update the JSON Object.",
					},
					"body": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Custom request body for the request",
					},
					"filter": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Custom result filter for the request",
					},
					"keepid": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Flag to keep the initial resource id, even when the response contains a new one",
					},
				},
			},
			"delete": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "manipulate the behavior for deleting the object",
				Attributes: map[string]schema.Attribute{
					"method": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "HTTP Method to use for the request",
						PlanModifiers: []planmodifier.String{
							modifier.StringDefaultValue(types.StringValue("DELETE")),
						},
					},
					"url": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Request URL to delete the JSON Object.",
					},
					"body": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Custom request body for the request",
					},
					"filter": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Custom result filter for the request",
					},
					"keepid": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Flag to keep the initial resource id, even when the response contains a new one",
					},
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
	var custom types2.CustomCreateMethod

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	customDiag := data.Create.As(ctx, &custom, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: false,
	})

	if customDiag != nil {
		resp.Diagnostics.Append(customDiag...)
		return
	}

	response := r.sendRequest(ctx, data, custom)
	tflog.Info(ctx, "<< http/"+fmt.Sprint(response.StatusCode()))
	tflog.Info(ctx, "<< "+response.Body())
	data.StatusCode = types.Int64Value(int64(response.StatusCode()))
	if (data.IgnoreStatusCode.ValueBool()) || (response.StatusCode() == 200) {
		jsondoc, _ := jsonquery.Parse(strings.NewReader(response.Body()))
		if !data.Filter.IsNull() {
			data.Result = types.StringValue(fmt.Sprint(jsonquery.FindOne(jsondoc, data.Filter.ValueString()).Value()))
		} else {
			data.Result = types.StringValue(response.Body())
		}
		data.ID = types.StringValue(fmt.Sprint(jsonquery.FindOne(jsondoc, data.Key.ValueString()).Value()))
	} else {
		resp.Diagnostics.AddError("Failed to create data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RestDataCall) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RestDataCallModel
	var custom types2.CustomCallAPI

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	customDiag := data.Read.As(ctx, &custom, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: false,
	})

	if customDiag != nil {
		resp.Diagnostics.Append(customDiag...)
		return
	}
	response := r.sendRequest(ctx, data, custom)
	tflog.Info(ctx, "<< http/"+fmt.Sprint(response.StatusCode()))
	tflog.Info(ctx, "<< "+response.Body())
	if (data.IgnoreStatusCode.ValueBool()) || (response.StatusCode() == 200) {
		jsondoc, _ := jsonquery.Parse(strings.NewReader(response.Body()))
		if !data.Filter.IsNull() || !custom.Filter.IsNull() {
			filter := data.Filter
			if !custom.Filter.IsNull() {
				filter = custom.Filter
			}
			data.Result = types.StringValue(fmt.Sprint(jsonquery.FindOne(jsondoc, filter.ValueString()).Value()))
		} else {
			data.Result = types.StringValue(response.Body())
		}
		if !custom.KeepId.ValueBool() {
			data.ID = types.StringValue(fmt.Sprint(jsonquery.FindOne(jsondoc, data.Key.ValueString()).Value()))
		}
	} else {
		resp.Diagnostics.AddError("Failed to fetch data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RestDataCall) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RestDataCallModel
	var custom types2.CustomCallAPI

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	customDiag := data.Update.As(ctx, &custom, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: false,
	})

	if customDiag != nil {
		resp.Diagnostics.Append(customDiag...)
		return
	}
	response := r.sendRequest(ctx, data, custom)
	tflog.Info(ctx, "<< http/"+fmt.Sprint(response.StatusCode()))
	tflog.Info(ctx, "<< "+response.Body())
	data.StatusCode = types.Int64Value(int64(response.StatusCode()))
	if (data.IgnoreStatusCode.ValueBool()) || (response.StatusCode() == 200) {
		jsondoc, _ := jsonquery.Parse(strings.NewReader(response.Body()))

		if !data.Filter.IsNull() || !custom.Filter.IsNull() {
			filter := data.Filter
			if !custom.Filter.IsNull() {
				filter = custom.Filter
			}
			data.Result = types.StringValue(fmt.Sprint(jsonquery.FindOne(jsondoc, filter.ValueString()).Value()))
		} else {
			data.Result = types.StringValue(response.Body())
		}
		if !custom.KeepId.ValueBool() {
			data.ID = types.StringValue(fmt.Sprint(jsonquery.FindOne(jsondoc, data.Key.ValueString()).Value()))
		}
	} else {
		resp.Diagnostics.AddError("Failed to update data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RestDataCall) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RestDataCallModel
	var custom types2.CustomCallAPI

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	customDiag := data.Delete.As(ctx, &custom, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: false,
	})

	if customDiag != nil {
		resp.Diagnostics.Append(customDiag...)
		return
	}
	response := r.sendRequest(ctx, data, custom)
	tflog.Info(ctx, "<< http/"+fmt.Sprint(response.StatusCode()))
	tflog.Info(ctx, "<< "+response.Body())
	if (data.IgnoreStatusCode.ValueBool()) || (response.StatusCode() == 200) {
		data.Result = types.StringValue(response.Body())
	} else {
		resp.Diagnostics.AddError("Failed to fetch data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
}
