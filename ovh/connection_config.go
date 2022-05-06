package ovh

import (
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/schema"
)

type ovhConfig struct {
	ApplicationKey    *string `cty:"application_key"`
	ApplicationSecret *string `cty:"application_secret"`
	ConsumerKey       *string `cty:"consumer_key"`
	Endpoint          *string `cty:"endpoint"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"application_key": {
		Type: schema.TypeString,
	},
	"application_secret": {
		Type: schema.TypeString,
	},
	"consumer_key": {
		Type: schema.TypeString,
	},
	"endpoint": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &ovhConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) ovhConfig {
	if connection == nil || connection.Config == nil {
		return ovhConfig{}
	}
	config, _ := connection.Config.(ovhConfig)
	return config
}
