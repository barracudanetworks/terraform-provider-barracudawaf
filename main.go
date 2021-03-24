package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hhakkaev/terraform-provider-barracudawaf/barracudawaf"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: barracudawaf.Provider,
	})
}
