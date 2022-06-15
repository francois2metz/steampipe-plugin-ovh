package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableOvhCloudVolume() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_volume",
		Description: "A volume is an independent additional disk.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listVolume,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getVolume,
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
				Description: "Volume ID.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Volume name.",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Description: "Volume region.",
			},
			{
				Name:        "attached_to",
				Type:        proto.ColumnType_JSON,
				Description: "Volume attached to instances ID.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Volume creation date.",
				Transform:   transform.FromField("CreationDate"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "Volume description.",
			},
			{
				Name:        "size",
				Type:        proto.ColumnType_INT,
				Description: "Volume size (in GB).",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Volume status.",
			},
			{
				Name:        "bootable",
				Type:        proto.ColumnType_BOOL,
				Description: "Volume bootable.",
			},
			{
				Name:        "planCode",
				Type:        proto.ColumnType_STRING,
				Description: "Order plan code.",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Volume type (classic, high-speed, high-speed-gen2",
			},
		},
	}
}

type Volume struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Region       string    `json:"region"`
	AttachedTo   []string  `json:"attachedTo"`
	CreationDate time.Time `json:"creationDate"`
	Description  string    `json:"description"`
	Size         int       `json:"size"`
	Status       string    `json:"status"`
	Bootable     bool      `json:"bootable"`
	PlanCode     string    `json:"planCode"`
	Type         string    `json:"type"`
}

func listVolume(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_volume.listVolume", "connection_error", err)
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	var volumes []Volume
	err = client.Get(fmt.Sprintf("/cloud/project/%s/volume", projectId), &volumes)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_volume.listVolume", err)
		return nil, err
	}
	for _, volume := range volumes {
		d.StreamListItem(ctx, volume)
	}
	return nil, nil
}

func getVolume(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_volume.getVolume", "connection_error", err)
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	id := d.KeyColumnQuals["id"].GetStringValue()
	var volume Volume
	err = client.Get(fmt.Sprintf("/cloud/project/%s/volume/%s", projectId, id), &volume)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_volume.getVolume", err)
		return nil, err
	}
	return volume, nil
}
