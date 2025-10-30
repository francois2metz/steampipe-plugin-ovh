package ovh

import (
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type ovhConfig struct {
	ApplicationKey    *string  `cty:"application_key"`
	ApplicationSecret *string  `cty:"application_secret"`
	ConsumerKey       *string  `cty:"consumer_key"`
	Endpoint          *string  `cty:"endpoint"`
	Regions           []string `cty:"regions"`
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
	"regions": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
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

	// Set default regions to ["*"] if not specified
	if config.Regions == nil {
		config.Regions = []string{"*"}
	}

	if len(config.Regions) == 0 {
		// Setting "regions = []" in the connection config is not valid
		errorMessage := fmt.Sprintf("connection %s has invalid value for \"regions\", it must contain at least 1 region.", connection.Name)
		panic(errorMessage)
	}

	return config
}
