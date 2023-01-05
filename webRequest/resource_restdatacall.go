package webRequest

import (
	"context"
	"curl-terraform-provider/client"
	"curl-terraform-provider/helper/modifier"
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

type CustomAPICall struct {
	Method types.String `tfsdk:"method"`
	URL    types.String `tfsdk:"url"`
	Body   types.String `tfsdk:"body"`
}

func NewRestDataCall() resource.Resource {
	return &RestDataCall{}
}

func (r *RestDataCall) sendRequest(ctx context.Context, data RestDataCallModel, custom CustomAPICall) client.Response {
	if custom.URL.IsNull() {
		custom.URL = types.StringValue(data.URL.ValueString() + "/" + data.ID.ValueString())
	}
	if custom.Body.IsNull() {
		custom.Body = data.Body
	}

	regex := regexp.MustCompile("{ID}")

	tflog.Info(ctx, ">> "+custom.Method.ValueString()+" "+regex.ReplaceAllString(custom.URL.ValueString(), data.ID.ValueString()))
	tflog.Info(ctx, ">> "+regex.ReplaceAllString(custom.Body.ValueString(), data.ID.ValueString()))
	request := r.client.NewRequest().SetMethod(custom.Method.ValueString()).SetURL(regex.ReplaceAllString(custom.URL.ValueString(), data.ID.ValueString())).SetBody(regex.ReplaceAllString(custom.Body.ValueString(), data.ID.ValueString()))

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
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("url"),
			path.MatchRoot("create"),
		),
	}
}

func (r *RestDataCall) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource interact with any rest like endpoint. All CRUD types are handled during lifetime of " +
			"this resource. fresh resource generate a create request, refreshing resource state result in a read action, " +
			"updating a attribute perform after apply a update action (partial update isn't supported) and finally deleting the resource performs a delete request " +
			"the primary key ob the resulting object is append to the url",

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
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("create"),
					),
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
				MarkdownDescription: "primary key of the received object, to generate/manipulate the request url",
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
			"create": schema.ObjectAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "manipulate the behavior for object creation",
				AttributeTypes: map[string]attr.Type{
					"method": types.StringType,
					"url":    types.StringType,
					"body":   types.StringType,
				},
				PlanModifiers: []planmodifier.Object{
					modifier.CustomCallDefaultValueModifier("POST"),
				},
			},
			"read": schema.ObjectAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "manipulate the behavior for reading the object",
				AttributeTypes: map[string]attr.Type{
					"method": types.StringType,
					"url":    types.StringType,
					"body":   types.StringType,
				},
				PlanModifiers: []planmodifier.Object{
					modifier.CustomCallDefaultValueModifier("GET"),
				},
			},
			"update": schema.ObjectAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "manipulate the behavior for updating the object",
				AttributeTypes: map[string]attr.Type{
					"method": types.StringType,
					"url":    types.StringType,
					"body":   types.StringType,
				}, PlanModifiers: []planmodifier.Object{
					modifier.CustomCallDefaultValueModifier("PUT"),
				},
			},
			"delete": schema.ObjectAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "manipulate the behavior for deleting the object",
				AttributeTypes: map[string]attr.Type{
					"method": types.StringType,
					"url":    types.StringType,
					"body":   types.StringType,
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

	response := r.sendRequest(ctx, data, custom)
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
	response := r.sendRequest(ctx, data, custom)
	if (data.IgnoreStatusCode.ValueBool()) || (response.StatusCode() == 200) {
		data.Result = types.StringValue(response.Body())
	} else {
		resp.Diagnostics.AddError("Failed to fetch data", "received response code "+fmt.Sprint(response.StatusCode()))
	}
}
