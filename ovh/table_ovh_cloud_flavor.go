package ovh

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudFlavor() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_flavor",
		Description: "A flavor is the instance model defining its characteristics in terms of resources.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listFlavor,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getFlavor,
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
				Description: "Flavor ID.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Flavor name.",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Description: "Flavor region.",
			},
			{
				Name:        "ram",
				Type:        proto.ColumnType_INT,
				Description: "Ram quantity (Gio).",
				Transform:   transform.FromField("Ram"),
			},
			{
				Name:        "disk",
				Type:        proto.ColumnType_INT,
				Description: "Number of disk.",
			},
			{
				Name:        "vcpus",
				Type:        proto.ColumnType_INT,
				Description: "Number of VCPUs.",
				Transform:   transform.FromField("VCPUs"),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Flavor type.",
			},
			{
				Name:        "os_type",
				Type:        proto.ColumnType_STRING,
				Description: "OS to install on.",
				Transform:   transform.FromField("OSType"),
			},
			{
				Name:        "inbound_bandwidth",
				Type:        proto.ColumnType_INT,
				Description: "Max capacity of inbound traffic in Mbit/s.",
			},
			{
				Name:        "outbound_bandwidth",
				Type:        proto.ColumnType_INT,
				Description: "Max capacity of outbound traffic in Mbit/s.",
			},
			{
				Name:        "available",
				Type:        proto.ColumnType_BOOL,
				Description: "Available in stock.",
			},
			{
				Name:        "quota",
				Type:        proto.ColumnType_INT,
				Description: "Number instance you can spawn with your actual quota.",
			},
			{
				Name:        "plan_codes_monthly",
				Type:        proto.ColumnType_STRING,
				Description: "Plan code to order monthly instance",
				Transform:   transform.FromField("PlanCodes.Monthly"),
			},
			{
				Name:        "plan_codes_hourly",
				Type:        proto.ColumnType_STRING,
				Description: "Plan code to order hourly instance",
				Transform:   transform.FromField("PlanCodes.Hourly"),
			},
		},
	}
}

type PlanCodes struct {
	Monthly string `string:"monthly"`
	Hourly  string `string:"hourly"`
}

type Flavor struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Region            string    `json:"region"`
	Ram               int       `json:"ram"`
	Disk              int       `json:"disk"`
	VCPUs             int       `json:"vcpus"`
	Type              string    `json:"type"`
	OSType            string    `json:"osType"`
	InboundBandwidth  int       `json:"inboundBandwidth"`
	OutboundBandwidth int       `json:"outboundBandwidth"`
	Available         bool      `json:"available"`
	Quota             int       `json:"quota"`
	PlanCodes         PlanCodes `json:"planCodes"`
}

func listFlavor(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_flavor.listFlavor", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	var flavors []Flavor
	err = client.Get(fmt.Sprintf("/cloud/project/%s/flavor", projectId), &flavors)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_flavor.listFlavor", err)
		return nil, err
	}
	for _, flavor := range flavors {
		d.StreamListItem(ctx, flavor)
	}
	return nil, nil
}

func getFlavor(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_flavor.getFlavor", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()
	var flavor Flavor
	err = client.Get(fmt.Sprintf("/cloud/project/%s/flavor/%s", projectId, id), &flavor)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_flavor.getFlavor", err)
		return nil, err
	}
	return flavor, nil
}
