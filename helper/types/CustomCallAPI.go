package types

import "github.com/hashicorp/terraform-plugin-framework/types"

var _ AlternativeRequestParameter = &CustomCallAPI{}

type CustomCallAPI struct {
	Method types.String `tfsdk:"method"`
	URL    types.String `tfsdk:"url"`
	Body   types.String `tfsdk:"body"`
	KeepId types.Bool   `tfsdk:"keepid"`
	Filter types.String `tfsdk:"filter"`
}

func (c CustomCallAPI) GetMethod() types.String {
	return c.Method
}

func (c CustomCallAPI) GetURL() types.String {
	return c.URL
}

func (c CustomCallAPI) GetBody() types.String {
	return c.Body
}
