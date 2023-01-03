//go:build tools

package tools

import (
	_ "github.com/goreleaser/goreleaser"
	// Documentation generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
