package ovh

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudRegion() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_region",
		Description: "Cloud regions.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listRegion,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "name"}),
			Hydrate:    getRegion,
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
				Description: "Name of the region.",
			},
			{
				Name:        "continent_code",
				Hydrate:     getRegionInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Region continent code.",
			},
			{
				Name:        "datacenter_location",
				Hydrate:     getRegionInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Location of the datacenter where the region is.",
			},
			{
				Name:        "ip_countries",
				Hydrate:     getRegionInfo,
				Type:        proto.ColumnType_JSON,
				Description: "Allowed countries for failover ip.",
			},
			{
				Name:        "services",
				Hydrate:     getRegionInfo,
				Type:        proto.ColumnType_JSON,
				Description: "Details about components status.",
			},
			{
				Name:        "status",
				Hydrate:     getRegionInfo,
				Type:        proto.ColumnType_STRING,
				Description: "OpenStack region status.",
			},
			{
				Name:        "type",
				Hydrate:     getRegionInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Region type.",
			},
		},
	}

}

type Region struct {
	Name               string      `json:"name"`
	ContinentCode      string      `json:"continentCode"`
	DatacenterLocation string      `json:"datacenterLocation"`
	IpCountries        []string    `json:"ipCountries"`
	Services           []Component `json:"services"`
	Status             string      `json:"status"`
	Type               string      `json:"type"`
}

type Component struct {
	Endpoint string `json:"endpoint"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

func getRegionInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := h.Item.(Region)
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_region.getRegionInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/cloud/project/%s/region/%s", projectId, region.Name), &region)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_region.getRegionInfo", err)
		return nil, err
	}
	return region, nil
}

func listRegion(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_region.listRegion", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	var regionNames []string
	err = client.Get(fmt.Sprintf("/cloud/project/%s/region", projectId), &regionNames)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_region.listRegion", err)
		return nil, err
	}
	for _, regionName := range regionNames {
		var region Region
		region.Name = regionName
		d.StreamListItem(ctx, region)
	}
	return nil, nil
}

func getRegion(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	var region Region
	region.Name = name
	return region, nil
}
