package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ovh/go-ovh/ovh"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableOvhCloudLoadBalancer() *plugin.Table {
	return &plugin.Table{
		Name:              "ovh_cloud_loadbalancer",
		Description:       "A load balancer distributes incoming traffic across a pool of backend servers.",
		GetMatrixItemFunc: RegionMatrix("octavialoadbalancer"),
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "project_id", Require: plugin.Required},
			},
			Hydrate: listLoadBalancer,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getLoadBalancer,
			// Every region of the project offering the service is queried in turn, so a get by
			// id necessarily misses in all but the one region holding the load balancer. Without
			// this the misses fail the whole query instead of returning the single matching row.
			// Deliberately not set on the list: there a 404 would mean the endpoint moved, and
			// silently returning no load balancer at all is worse than failing loudly.
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreNotFound,
			},
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
				Description: "ID of the load balancer.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the load balancer.",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Description: "Region of the load balancer.",
			},
			{
				Name:        "flavor_id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the flavor of the load balancer.",
			},
			{
				Name:        "operating_status",
				Type:        proto.ColumnType_STRING,
				Description: "Operating status of the load balancer: degraded, draining, error, noMonitor, offline or online.",
			},
			{
				Name:        "provisioning_status",
				Type:        proto.ColumnType_STRING,
				Description: "Provisioning status of the load balancer: active, creating, deleted, deleting, error or updating.",
			},
			{
				Name:        "vip_address",
				Type:        proto.ColumnType_STRING,
				Description: "IP address of the virtual IP.",
			},
			{
				Name:        "vip_network_id",
				Type:        proto.ColumnType_STRING,
				Description: "Openstack ID of the network of the virtual IP.",
			},
			{
				Name:        "vip_subnet_id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the subnet of the virtual IP.",
			},
			{
				Name:        "floating_ip_id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the floating IP associated with the load balancer, if any.",
				Transform:   transform.FromField("FloatingIP.ID").NullIfZero(),
			},
			{
				Name:        "floating_ip",
				Type:        proto.ColumnType_STRING,
				Description: "Floating IP address associated with the load balancer, if any.",
				Transform:   transform.FromField("FloatingIP.IP").NullIfZero(),
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and timestamp when the load balancer was created.",
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and timestamp when the load balancer was last updated.",
			},
		},
	}
}

type LoadBalancer struct {
	CreatedAt          time.Time              `json:"createdAt"`
	FlavorID           string                 `json:"flavorId"`
	FloatingIP         LoadBalancerFloatingIP `json:"floatingIp"`
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	OperatingStatus    string                 `json:"operatingStatus"`
	ProvisioningStatus string                 `json:"provisioningStatus"`
	Region             string                 `json:"region"`
	UpdatedAt          time.Time              `json:"updatedAt"`
	VipAddress         string                 `json:"vipAddress"`
	VipNetworkID       string                 `json:"vipNetworkId"`
	VipSubnetID        string                 `json:"vipSubnetId"`
}

type LoadBalancerFloatingIP struct {
	ID string `json:"id"`
	IP string `json:"ip"`
}

// shouldIgnoreNotFound reports whether the error is an OVHcloud API 404.
func shouldIgnoreNotFound(_ context.Context, _ *plugin.QueryData, _ *plugin.HydrateData, err error) bool {
	var apiError *ovh.APIError
	if errors.As(err, &apiError) {
		return apiError.Code == http.StatusNotFound
	}
	return false
}

func listLoadBalancer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_loadbalancer.listLoadBalancer", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	region := d.EqualsQualString("region")

	var loadBalancers []LoadBalancer
	err = client.Get(fmt.Sprintf("/cloud/project/%s/region/%s/loadbalancing/loadbalancer", projectId, region), &loadBalancers)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_loadbalancer.listLoadBalancer", "api_error", err, "region", region)
		return nil, err
	}
	for _, loadBalancer := range loadBalancers {
		d.StreamListItem(ctx, loadBalancer)
	}
	return nil, nil
}

func getLoadBalancer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_loadbalancer.getLoadBalancer", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	region := d.EqualsQualString("region")
	id := d.EqualsQuals["id"].GetStringValue()

	var loadBalancer LoadBalancer
	err = client.Get(fmt.Sprintf("/cloud/project/%s/region/%s/loadbalancing/loadbalancer/%s", projectId, region, id), &loadBalancer)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_loadbalancer.getLoadBalancer", "api_error", err, "region", region, "id", id)
		return nil, err
	}
	return loadBalancer, nil
}
