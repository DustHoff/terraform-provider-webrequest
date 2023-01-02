package webRequest

import (
	"context"
	client2 "curl-terraform-provider/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &WebRequestProvider{}

type WebRequestProvider struct {
	// Version is an example field that can be set with an actual provider
	// version on release, "dev" when the provider is built and ran locally,
	// and "test" when running acceptance testing.
	version string
}

type WebRequestProviderModel struct {
	Timeout types.Int64 `tfsdk:"timeout"`
}

func NewProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &WebRequestProvider{
			version: version,
		}
	}
}
func (p *WebRequestProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "webrequest"
}

func (p *WebRequestProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This provider can interact with any rest compatible endpoints. It work like the command line tool" +
			"_curl_",
		Attributes: map[string]schema.Attribute{
			"timeout": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Request Timeout",
			},
		},
	}
}

func (p *WebRequestProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data WebRequestProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Timeout.IsNull() || data.Timeout.IsUnknown() {
		data.Timeout = types.Int64Value(60)
	}
	client := client2.NewClient(int(data.Timeout.ValueInt64()))
	resp.ResourceData = client
	resp.DataSourceData = client
}
func (p *WebRequestProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWebRequestDataSource,
	}
}

func (p *WebRequestProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRestDataCall,
	}
}
