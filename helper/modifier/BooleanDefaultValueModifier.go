package modifier

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.Bool = (*BooleanDefaultValuePlanModifier)(nil)

type BooleanDefaultValuePlanModifier struct {
	DefaultValue types.Bool
}

func BooleanDefaultValueModifier(value types.Bool) *BooleanDefaultValuePlanModifier {
	return &BooleanDefaultValuePlanModifier{
		DefaultValue: value,
	}
}

func (b BooleanDefaultValuePlanModifier) Description(ctx context.Context) string {
	return ""
}

func (b BooleanDefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (b BooleanDefaultValuePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, res *planmodifier.BoolResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	res.PlanValue = b.DefaultValue
}
