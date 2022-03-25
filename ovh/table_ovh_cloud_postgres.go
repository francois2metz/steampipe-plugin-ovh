package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableOvhCloudPostgres() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_postgres",
		Description: "List all the postgresql of the project.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listPostgres,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getPostgres,
		},
		HydrateDependencies: []plugin.HydrateDependencies{
			{Func: getProjectInfo},
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
				Description: "Service id.",
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
				Type:        proto.ColumnType_DATETIME,
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

type MaintenanceWindow struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Postgres struct {
	ID                string            `json:"id"`
	CreatedAt         *time.Time        `json:"createdAt"`
	Plan              string            `json:"plan"`
	Engine            string            `json:"engine"`
	Status            string            `json:"status"`
	NodeNumber        int               `json:"nodeNumber"`
	Description       string            `json:"description"`
	MaintenanceWindow MaintenanceWindow `json:"maintenance_window"`
	Version           string            `json:"version"`
	NetworkType       string            `json:"networkType"`
	Flavor            string            `json:"flavor"`
}

func getPostgresInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPostgresInfo")
	postgres := h.Item.(Postgres)
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/cloud/project/%s/database/postgresql/%s", projectId, postgres.ID), &postgres)
	if err != nil {
		return nil, err
	}
	return postgres, nil
}

func listPostgres(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	var postgresIds []string
	err = client.Get(fmt.Sprintf("/cloud/project/%s/database/postgresql", projectId), &postgresIds)
	if err != nil {
		return nil, err
	}
	for _, postgresId := range postgresIds {
		var postgres Postgres
		postgres.ID = postgresId
		d.StreamListItem(ctx, postgres)
	}
	return nil, nil
}

func getPostgres(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.KeyColumnQuals["id"].GetStringValue()
	var postgres Postgres
	postgres.ID = id
	return postgres, nil
}
