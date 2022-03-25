package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableOvhCloudImage() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_image",
		Description: "Get images.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listImage,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getImage,
		},
		Columns: []*plugin.Column{
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("project_id"),
				Description: "Project id.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Image id.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Image name.",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Description: "Image region.",
			},
			{
				Name:        "visibility",
				Type:        proto.ColumnType_STRING,
				Description: "Image visibility.",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Image type.",
			},
			{
				Name:        "min_disk",
				Type:        proto.ColumnType_INT,
				Description: "Minimum disks required to use image.",
				Transform:   transform.FromField("MinDisk"),
			},
			{
				Name:        "min_ram",
				Type:        proto.ColumnType_INT,
				Description: "Minimum RAM required to use image.",
				Transform:   transform.FromField("MinRam"),
			},
			{
				Name:        "size",
				Type:        proto.ColumnType_DOUBLE,
				Description: "Image size (in GiB).",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Image creation date.",
				Transform:   transform.FromField("CreationDate"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Image status.",
			},
			{
				Name:        "user",
				Type:        proto.ColumnType_STRING,
				Description: "User to connect with.",
			},
			{
				Name:        "flavor_type",
				Type:        proto.ColumnType_STRING,
				Description: "Image usable only for this type of flavor if not null.",
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: "Tags about the image.",
			},
			{
				Name:        "plan_code",
				Type:        proto.ColumnType_STRING,
				Description: "Order plan code.",
			},
		},
	}
}

type Image struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Region       string    `json:"region"`
	Visibility   string    `json:"visibility"`
	Type         string    `json:"type"`
	MinDisk      int       `json:"minDisk"`
	MinRam       int       `json:"minRam"`
	Size         float32   `json:"size"`
	CreationDate time.Time `json:"creationDate"`
	Status       string    `json:"status"`
	User         string    `json:"user"`
	FlavorType   string    `json:"flavorType"`
	Tags         []string  `json:"tags"`
	PlanCode     string    `json:"planCode"`
}

func listImage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	var images []Image
	err = client.Get(fmt.Sprintf("/cloud/project/%s/image", projectId), &images)
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		d.StreamListItem(ctx, image)
	}
	return nil, nil
}

func getImage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	id := d.KeyColumnQuals["id"].GetStringValue()
	var image Image
	err = client.Get(fmt.Sprintf("/cloud/project/%s/image/%s", projectId, id), &image)
	if err != nil {
		return nil, err
	}
	return image, nil
}
