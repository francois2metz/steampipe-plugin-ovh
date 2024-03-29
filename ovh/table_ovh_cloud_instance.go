package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudInstance() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_instance",
		Description: "An instance is a virtual server in the OVH cloud.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listInstance,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getInstance,
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
				Description: "Instance ID.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Instance name.",
			},
			{
				Name:        "flavor_id",
				Type:        proto.ColumnType_STRING,
				Description: "Instance flavor ID.",
			},
			{
				Name:        "image_id",
				Type:        proto.ColumnType_STRING,
				Description: "Instance image ID.",
			},
			{
				Name:        "ssh_key_id",
				Type:        proto.ColumnType_STRING,
				Description: "Instance ssh key ID.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Instance creation date.",
				Transform:   transform.FromField("Created"),
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Description: "Region of the instance.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Instance status.",
			},
			{
				Name:        "plan_code",
				Type:        proto.ColumnType_STRING,
				Description: "Order plan code.",
			},
			{
				Name:        "current_month_outgoing_traffic",
				Type:        proto.ColumnType_INT,
				Description: "Instance outgoing network traffic for the current month (in bytes).",
			},
		},
	}
}

type Instance struct {
	ID                          string    `json:"id"`
	Name                        string    `json:"name"`
	FlavorID                    string    `json:"flavorId"`
	Flavor                      Flavor    `json:"flavor"`
	ImageID                     string    `json:"imageId"`
	Image                       Image     `json:"image"`
	SSHKeyID                    string    `json:"sshKeyId"`
	SSHKey                      SshKey    `json:"sshKey"`
	Created                     time.Time `json:"created"`
	Region                      string    `json:"region"`
	Status                      string    `json:"status"`
	PlanCode                    string    `json:"planCode"`
	CurrentMonthOutgoingTraffic *int      `json:"currentMonthOutgoingTraffic,omitempty"`
}

func listInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_instance.listInstance", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	var instances []Instance
	err = client.Get(fmt.Sprintf("/cloud/project/%s/instance", projectId), &instances)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_instance.listInstance", err)
		return nil, err
	}
	for _, instance := range instances {
		d.StreamListItem(ctx, instance)
	}
	return nil, nil
}

func getInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_instance.getInstance", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()
	var instance Instance
	err = client.Get(fmt.Sprintf("/cloud/project/%s/instance/%s", projectId, id), &instance)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_instance.getInstance", err)
		return nil, err
	}
	instance.ImageID = instance.Image.ID
	instance.FlavorID = instance.Flavor.ID
	instance.SSHKeyID = instance.SSHKey.ID
	return instance, nil
}
