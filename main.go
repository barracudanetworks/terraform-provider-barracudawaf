package main

import (
	"github.com/hhakkaev/terraform-provider-barracudawaf/barracudawaf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: barracudawaf.Provider,
	})
}
