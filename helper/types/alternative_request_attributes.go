package types

import "github.com/hashicorp/terraform-plugin-framework/types"

type AlternativeRequestParameter interface {
	GetMethod() types.String
	GetURL() types.String
	GetBody() types.String
}
