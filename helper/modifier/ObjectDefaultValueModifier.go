package modifier

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.Object = (*ObjectDefaultValuePlanModifier)(nil)

type ObjectDefaultValuePlanModifier struct {
	DefaultValue types.Object
}

func CustomCallDefaultValueModifier(method string) *ObjectDefaultValuePlanModifier {
	objectValue, _ := types.ObjectValue(map[string]attr.Type{
		"method": types.StringType,
		"url":    types.StringType,
		"body":   types.StringType,
	}, map[string]attr.Value{
		"method": types.StringValue(method),
		"url":    types.StringNull(),
		"body":   types.StringNull(),
	})
	return &ObjectDefaultValuePlanModifier{
		DefaultValue: objectValue,
	}
}

func (o ObjectDefaultValuePlanModifier) Description(ctx context.Context) string {
	return ""
}

func (o ObjectDefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (o ObjectDefaultValuePlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, res *planmodifier.ObjectResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	res.PlanValue = o.DefaultValue
}
