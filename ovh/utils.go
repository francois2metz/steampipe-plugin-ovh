package ovh

import (
	"context"
	"errors"

	"github.com/ovh/go-ovh/ovh"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (*ovh.Client, error) {
	// get ovh client from cache
	cacheKey := "ovh"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*ovh.Client), nil
	}

	applicationKey := ""
	applicationSecret := ""
	consumerKey := ""
	endpoint := ""

	ovhConfig := GetConfig(d.Connection)
	if &ovhConfig != nil {
		if ovhConfig.ApplicationKey != nil {
			applicationKey = *ovhConfig.ApplicationKey
		}
		if ovhConfig.ApplicationSecret != nil {
			applicationSecret = *ovhConfig.ApplicationSecret
		}
		if ovhConfig.ConsumerKey != nil {
			consumerKey = *ovhConfig.ConsumerKey
		}
		if ovhConfig.Endpoint != nil {
			endpoint = *ovhConfig.Endpoint
		}
	}

	if applicationKey == "" {
		return nil, errors.New("'application_key' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if applicationSecret == "" {
		return nil, errors.New("'application_secret' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if consumerKey == "" {
		return nil, errors.New("'consumer_key' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if endpoint == "" {
		return nil, errors.New("'endpoint' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	client, _ := ovh.NewClient(
		endpoint,
		applicationKey,
		applicationSecret,
		consumerKey,
	)

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}
