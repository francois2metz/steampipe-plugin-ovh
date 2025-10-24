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
			"ovh_bill":                  tableOvhBill(),
			"ovh_bill_detail":           tableOvhBillDetails(),
			"ovh_ceph":       			 tableOvhCeph(),
			"ovh_cloud_ai_app":          tableOvhCloudAIApp(),
			"ovh_cloud_ai_job":          tableOvhCloudAIJob(),
			"ovh_cloud_ai_notebook":     tableOvhCloudAINotebook(),
			"ovh_cloud_data_job":        tableOvhCloudDataJob(),
			"ovh_cloud_database":        tableOvhCloudDatabase(),
			"ovh_cloud_flavor":          tableOvhCloudFlavor(),
			"ovh_cloud_image":           tableOvhCloudImage(),
			"ovh_cloud_instance":        tableOvhCloudInstance(),
			"ovh_cloud_postgres":        tableOvhCloudPostgres(),
			"ovh_cloud_project":         tableOvhCloudProject(),
			"ovh_cloud_region":          tableOvhCloudRegion(),
			"ovh_cloud_ssh_key":         tableOvhCloudSshKey(),
			"ovh_cloud_storage_s3":      tableOvhCloudStorageS3(),
			"ovh_cloud_storage_swift":   tableOvhCloudStorageSwift(),
			"ovh_cloud_volume":          tableOvhCloudVolume(),
			"ovh_cloud_volume_snapshot": tableOvhCloudVolumeSnapshot(),
			"ovh_iam_resource":          tableOvhIamResource(),
			"ovh_log_self":              tableOvhLog(),
			"ovh_refund":                tableOvhRefund(),
			"ovh_refund_detail":         tableOvhRefundDetails(),
		},
	}
	return p
}
