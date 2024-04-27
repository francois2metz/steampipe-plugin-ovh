package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudS3Storage() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_s3_storage",
		Description: "A S3 storage is an object storage.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "region"}),
			Hydrate:    listS3StorageContainer,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "region", "name"}),
			Hydrate:    getS3StorageContainer,
		},
		Columns: []*plugin.Column{
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("project_id"),
				Description: "Project ID.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Container name.",
			},
			{
				Name:        "virtual_host",
				Type:        proto.ColumnType_STRING,
				Description: "Container virtual host.",
			},
			{
				Name:        "owner_id",
				Type:        proto.ColumnType_INT,
				Description: "Container owner userID.",
			},
			{
				Name:        "objects_count",
				Type:        proto.ColumnType_INT,
				Description: "Container total objects count.",
			},
			{
				Name:        "objects_size",
				Type:        proto.ColumnType_INT,
				Description: "Container total objects size (bytes).",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Description: "Region of the container.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and timestamp when the resource was created.",
			},
			{
				Name:        "encryption_sse_algorithm",
				Type:        proto.ColumnType_STRING,
				Description: "Encryption configuration.",
				Transform:   transform.FromField("Encryption.SSEAlgorithm"),
			},
		},
	}
}

type S3StorageContainer struct {
	Name         string                       `json:"name"`
	VirtualHost  string                       `json:"virtualHost"`
	OwnerID      int                          `json:"ownerId"`
	ObjectsCount int                          `json:"objectsCount"`
	ObjectsSize  int                          `json:"objectsSize"`
	Region       string                       `json:"region"`
	CreatedAt    time.Time                    `json:"createdAt"`
	Encryption   S3StorageContainerEncryption `json:"encryption"`
}
type S3StorageContainerEncryption struct {
	SSEAlgorithm string `json:"sseAlgorithm"`
}

func listS3StorageContainer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_s3_storage.listS3StorageContainer", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	region := d.EqualsQuals["region"].GetStringValue()

	var containers []S3StorageContainer
	err = client.Get(fmt.Sprintf("/cloud/project/%s/region/%s/storage", projectId, region), &containers)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_s3_storage.listS3StorageContainer", err)
		return nil, err
	}
	for _, container := range containers {
		d.StreamListItem(ctx, container)
	}
	return nil, nil
}

func getS3StorageContainer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_s3_storage.getS3StorageContainer", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	region := d.EqualsQuals["region"].GetStringValue()
	name := d.EqualsQuals["name"].GetStringValue()
	var container S3StorageContainer
	err = client.Get(fmt.Sprintf("/cloud/project/%s/region/%s/storage/%s", projectId, region, name), &container)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_s3_storage.getS3StorageContainer", err)
		return nil, err
	}
	return container, nil
}
