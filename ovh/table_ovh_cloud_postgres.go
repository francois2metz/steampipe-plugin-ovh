package ovh

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudPostgres() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_postgres",
		Description: "An hosted PostgreSQL database.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listPostgres,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getPostgres,
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
				Description: "Service ID.",
			},
			{
				Name:        "engine",
				Hydrate:     getPostgresInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Name of the engine of the service.",
			},
			{
				Name:        "plan",
				Hydrate:     getPostgresInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Plan of the cluster.",
			},
			{
				Name:        "created_at",
				Hydrate:     getPostgresInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date of the creation of the cluster.",
			},
			{
				Name:        "status",
				Hydrate:     getPostgresInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Current status of the cluster.",
			},
			{
				Name:        "node_number",
				Hydrate:     getPostgresInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Number of nodes in the cluster.",
			},
			{
				Name:        "description",
				Hydrate:     getPostgresInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Description of the cluster.",
			},
			{
				Name:        "version",
				Hydrate:     getPostgresInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Version of the engine deployed on the cluster.",
			},
			{
				Name:        "network_type",
				Hydrate:     getPostgresInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Type of network of the cluster.",
			},
			{
				Name:        "flavor",
				Hydrate:     getPostgresInfo,
				Type:        proto.ColumnType_STRING,
				Description: "The VM flavor used for this cluster.",
			},
		},
	}

}
func getPostgresInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	postgres := h.Item.(Database)
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_postgres.getPostgresInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/cloud/project/%s/database/postgresql/%s", projectId, postgres.ID), &postgres)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_postgres.getPostgresInfo", err)
		return nil, err
	}
	return postgres, nil
}

func listPostgres(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_postgres.listPostgres", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	var postgresIds []string
	err = client.Get(fmt.Sprintf("/cloud/project/%s/database/postgresql", projectId), &postgresIds)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_postgres.listPostgres", err)
		return nil, err
	}
	for _, postgresId := range postgresIds {
		var postgres Database
		postgres.ID = postgresId
		d.StreamListItem(ctx, postgres)
	}
	return nil, nil
}

func getPostgres(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetStringValue()
	var postgres Database
	postgres.ID = id
	return postgres, nil
}
