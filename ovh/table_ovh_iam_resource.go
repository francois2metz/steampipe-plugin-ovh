package ovh

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhIamResource() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_iam_resource",
		Description: "IAM resources in the OVH account.",
		List: &plugin.ListConfig{
			Hydrate: listIamResource,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Resource unique identifier (UUID).",
			},
			{
				Name:        "urn",
				Type:        proto.ColumnType_STRING,
				Description: "Unique resource name used in policies.",
				Transform: transform.FromField("URN"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Resource name.",
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Resource display name.",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Resource type.",
			},
			{
				Name:        "owner",
				Type:        proto.ColumnType_STRING,
				Description: "Resource owner.",
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: "Resource tags. Tags that were internally computed are prefixed with ovh:.",
			},
		},
	}
}

type IamResource struct {
	ID          string            `json:"id"`
	URN         string            `json:"urn"`
	Name        string            `json:"name"`
	DisplayName string            `json:"displayName,omitempty"`
	Type        string            `json:"type"`
	Owner       string            `json:"owner"`
	Tags        map[string]string `json:"tags,omitempty"`
}

func listIamResource(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_iam_resource.listIamResource", "connection_error", err)
		return nil, err
	}

	var resources []IamResource
	err = client.Get("/v2/iam/resource", &resources)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_iam_resource.listIamResource", err)
		return nil, err
	}

	for _, resource := range resources {
		d.StreamListItem(ctx, resource)
	}

	return nil, nil
}
