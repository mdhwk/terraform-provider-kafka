package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/mdhwk/terraform-provider-kafka/internal/provider"
)

var (
	version string = "dev"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: provider.New(version)}

	if debugMode {
		// TODO: update this string with the full name of your provider as used in your configs
		err := plugin.Debug(context.Background(), "registry.terraform.io/hashicorp/scaffolding", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
