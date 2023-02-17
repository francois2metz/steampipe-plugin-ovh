package main

import (
	"github.com/francois2metz/steampipe-plugin-ovh/ovh"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: ovh.Plugin})
}
