package main

import (
	"context"
	"curl-terraform-provider/webRequest"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name webrequest
var (
	// Example version string that can be overwritten by a release process
	version string = "1.0.0"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/DustHoff/webrequest",
	}

	err := providerserver.Serve(context.Background(), webRequest.NewProvider(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
