package ovh

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-ovh",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"ovh_cloud_instance": tableOvhCloudInstance(),
			"ovh_cloud_postgres": tableOvhCloudPostgres(),
			"ovh_cloud_project":  tableOvhCloudProject(),
			"ovh_cloud_sshkey":   tableOvhCloudSshkey(),
			"ovh_cloud_storage":  tableOvhCloudStorage(),
		},
	}
	return p
}
