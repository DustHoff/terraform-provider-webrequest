package types

import "github.com/hashicorp/terraform-plugin-framework/types"

var _ AlternativeRequestParameter = &CustomCreateMethod{}

type CustomCreateMethod struct {
	Method types.String `tfsdk:"method"`
}

func (c CustomCreateMethod) GetMethod() types.String {
	return c.Method
}

func (c CustomCreateMethod) GetURL() types.String {
	return types.StringNull()
}

func (c CustomCreateMethod) GetBody() types.String {
	return types.StringNull()
}
