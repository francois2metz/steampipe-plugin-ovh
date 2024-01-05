package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudVolumeSnapshot() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_volume_snapshot",
		Description: "A volume is an independent additional disk.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listVolumeSnapshot,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getVolumeSnapshot,
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
				Description: "ID.",
			},
			{
				Name:        "creationDate",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Volume creation date.",
				Transform:   transform.FromField("CreationDate"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Volume Snapshot Name.",
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "Volume Snapshot Description",
			},
			{
				Name:        "size",
				Type:        proto.ColumnType_INT,
				Description: "Volume Snapshot size (in GB).",
			},
			{
				Name:        "volumeId",
				Type:        proto.ColumnType_STRING,
				Description: "Volume Snapshot ID.",
				Transform:   transform.FromField("VolumeId"),
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Description: "Volume Snapshot Region.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Volume Snapshot Status. (available, creating, deleting, error, error_deleting)",
			},
			{
				Name:        "planCode",
				Type:        proto.ColumnType_STRING,
				Description: "Volume Snapshot Plan Code",
			},
		},
	}
}

type VolumeSnapShot struct {
	ID           string    `json:"id"`
	CreationDate time.Time `json:"creationDate"`
	Name	     string    `json:"name"`
	Description  string    `json:"description"`
	Size         int       `json:"size"`
	VolumeId     string    `json:"volumeId"`
	Region       string    `json:"region"`
	Status       string    `json:"status"`
	PlanCode     string    `json:"planCode"`
}

func listVolumeSnapshot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_volume.listVolumeSnapshot", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	var volumes []VolumeSnapShot
	err = client.Get(fmt.Sprintf("/cloud/project/%s/volume/snapshot", projectId), &volumes)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_volume.listVolumeSnapshot", err)
		return nil, err
	}
	for _, volume := range volumes {
		d.StreamListItem(ctx, volume)
	}
	return nil, nil
}

func getVolumeSnapshot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_volume.getVolumeSnapshot", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()
	var volume VolumeSnapShot
	err = client.Get(fmt.Sprintf("/cloud/project/%s/volume/snapshot/%s", projectId, id), &volume)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_volume.getVolumeSnapshot", err)
		return nil, err
	}
	return volume, nil
}
