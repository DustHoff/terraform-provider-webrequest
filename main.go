package main

import (
	"context"
	"curl-terraform-provider/webRequest"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// Example version string that can be overwritten by a release process
	version string = "dev"
)

func main() {
	opts := providerserver.ServeOpts{
		// TODO: Update this string with the published name of your provider.
		Address: "registry.terraform.io/example-namespace/example",
	}

	err := providerserver.Serve(context.Background(), webRequest.NewProvider(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
