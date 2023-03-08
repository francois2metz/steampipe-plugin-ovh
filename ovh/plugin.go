package ovh

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			"ovh_cloud_ai_app":      tableOvhCloudAIApp(),
			"ovh_cloud_ai_job":      tableOvhCloudAIJob(),
			"ovh_cloud_ai_notebook": tableOvhCloudAINotebook(),
			"ovh_cloud_data_job":    tableOvhCloudDataJob(),
			"ovh_cloud_database":    tableOvhCloudDatabase(),
			"ovh_cloud_flavor":      tableOvhCloudFlavor(),
			"ovh_cloud_image":       tableOvhCloudImage(),
			"ovh_cloud_instance":    tableOvhCloudInstance(),
			"ovh_cloud_postgres":    tableOvhCloudPostgres(),
			"ovh_cloud_project":     tableOvhCloudProject(),
			"ovh_cloud_ssh_key":     tableOvhCloudSshKey(),
			"ovh_cloud_storage":     tableOvhCloudStorage(),
			"ovh_cloud_volume":      tableOvhCloudVolume(),
			"ovh_bill":              tableOvhBill(),
		},
	}
	return p
}
