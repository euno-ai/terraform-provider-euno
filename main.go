package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/euno-ai/terraform-provider-euno/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// Remove any date and time prefix in the log output when the server is running
	// supported by the Terraform CLI. It does not affect log output when the server
	// is running in debug mode.
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/euno-ai/euno",
		Debug:   debug,
	})

	if err != nil {
		log.Printf("[ERROR] Failed to start provider server: %v", err)
	}
}
