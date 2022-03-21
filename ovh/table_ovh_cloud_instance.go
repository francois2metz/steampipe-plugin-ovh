package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableOvhCloudInstance() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_instance",
		Description: "Get cloud instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listInstance,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getInstance,
		},
		Columns: []*plugin.Column{
			{Name: "project_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("project_id"), Description: "Project id."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Instance id."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Instance name."},
			{Name: "flavor_id", Type: proto.ColumnType_STRING, Description: "Instance flavor id."},
			{Name: "image_id", Type: proto.ColumnType_STRING, Description: "Instance image id."},
			{Name: "ssh_key_id", Type: proto.ColumnType_STRING, Description: "Instance ssh key id."},
			{Name: "created_at", Type: proto.ColumnType_DATETIME, Description: "Instance creation date.", Transform: transform.FromField("Created")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "Region of the instance."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Instance status."},
			{Name: "plan_code", Type: proto.ColumnType_STRING, Description: "Order plan code."},
			{Name: "current_month_outgoing_traffic", Type: proto.ColumnType_INT, Description: "Instance outgoing network traffic for the current month (in bytes)."},
		},
	}
}

type Instance struct {
	ID                          string    `json:"id"`
	Name                        string    `json:"name"`
	FlavorID                    string    `json:"flavorId"`
	ImageID                     string    `json:"imageId"`
	SSHKeyID                    string    `json:"sshKeyId"`
	Created                     time.Time `json:"created"`
	Region                      string    `json:"region"`
	Status                      string    `json:"status"`
	PlanCode                    string    `json:"planCode"`
	CurrentMonthOutgoingTraffic *int      `json:"currentMonthOutgoingTraffic,omitempty"`
}

func listInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	var instances []Instance
	err = client.Get(fmt.Sprintf("/cloud/project/%s/instance", projectId), &instances)
	if err != nil {
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
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	id := d.KeyColumnQuals["id"].GetStringValue()
	var instance Instance
	err = client.Get(fmt.Sprintf("/cloud/project/%s/instance/%s", projectId, id), &instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
