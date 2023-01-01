package modifier

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func MapDefaultValue(v types.Map) planmodifier.Map {
	return &mapDefaultValuePlanModifier{v}
}
func EmptyMapDefaultValue() planmodifier.Map {
	var empty map[string]attr.Value
	mapValue, _ := types.MapValue(types.StringType, empty)
	return MapDefaultValue(mapValue)
}

type mapDefaultValuePlanModifier struct {
	DefaultValue types.Map
}

var _ planmodifier.Map = (*mapDefaultValuePlanModifier)(nil)

func (apm *mapDefaultValuePlanModifier) Description(ctx context.Context) string {
	return ""
}

func (apm *mapDefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (apm *mapDefaultValuePlanModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, res *planmodifier.MapResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	res.PlanValue = apm.DefaultValue
}
