package main

import (
	"github.com/autoscalr/terraform-provider-autoscalr/autoscalr"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: autoscalr.Provider})
}
