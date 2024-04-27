package ovh

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudSwiftStorage() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_swift_storage",
		Description: "A Swift storage is an object storage similar to S3.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listSwiftStorageContainer,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getSwiftStorageContainer,
		},
		Columns: []*plugin.Column{
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("project_id"),
				Description: "Project ID.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Container ID.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Container name.",
			},
			{
				Name:        "stored_objects",
				Type:        proto.ColumnType_INT,
				Description: "Total objects stored.",
			},
			{
				Name:        "stored_bytes",
				Type:        proto.ColumnType_INT,
				Description: "Total bytes stored.",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Description: "Region of the container.",
			},
		},
	}
}

type SwiftStorageContainer struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	StoredObjects int    `json:"storedObjects"`
	StoredBytes   int    `json:"storedBytes"`
	Region        string `json:"region"`
}

func listSwiftStorageContainer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_swift_storage.listSwiftStorageContainer", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	var containers []SwiftStorageContainer
	err = client.Get(fmt.Sprintf("/cloud/project/%s/storage", projectId), &containers)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_swift_storage.listSwiftStorageContainer", err)
		return nil, err
	}
	for _, container := range containers {
		d.StreamListItem(ctx, container)
	}
	return nil, nil
}

func getSwiftStorageContainer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_swift_storage.getSwiftStorageContainer", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()
	var container SwiftStorageContainer
	err = client.Get(fmt.Sprintf("/cloud/project/%s/storage/%s", projectId, id), &container)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_swift_storage.getSwiftStorageContainer", err)
		return nil, err
	}
	container.ID = id
	return container, nil
}
