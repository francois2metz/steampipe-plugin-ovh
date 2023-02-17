package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudDatabase() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_database",
		Description: "An hosted database service.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listDatabase,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getDatabase,
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
				Hydrate:     getDatabaseInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Name of the engine of the service.",
			},
			{
				Name:        "plan",
				Hydrate:     getDatabaseInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Plan of the cluster.",
			},
			{
				Name:        "created_at",
				Hydrate:     getDatabaseInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date of the creation of the cluster.",
			},
			{
				Name:        "status",
				Hydrate:     getDatabaseInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Current status of the cluster.",
			},
			{
				Name:        "node_number",
				Hydrate:     getDatabaseInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Number of nodes in the cluster.",
			},
			{
				Name:        "description",
				Hydrate:     getDatabaseInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Description of the cluster.",
			},
			{
				Name:        "version",
				Hydrate:     getDatabaseInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Version of the engine deployed on the cluster.",
			},
			{
				Name:        "network_type",
				Hydrate:     getDatabaseInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Type of network of the cluster.",
			},
			{
				Name:        "flavor",
				Hydrate:     getDatabaseInfo,
				Type:        proto.ColumnType_STRING,
				Description: "The VM flavor used for this cluster.",
			},
			{
				Name:        "backup_time",
				Hydrate:     getDatabaseInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Time on which backups start every day.",
			},
			{
				Name:        "maintenance_time",
				Hydrate:     getDatabaseInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Time on which maintenances can start every day.",
			},
		},
	}

}

type Database struct {
	ID              string     `json:"id"`
	CreatedAt       *time.Time `json:"createdAt"`
	Plan            string     `json:"plan"`
	Engine          string     `json:"engine"`
	Status          string     `json:"status"`
	NodeNumber      int        `json:"nodeNumber"`
	Description     string     `json:"description"`
	Version         string     `json:"version"`
	NetworkType     string     `json:"networkType"`
	Flavor          string     `json:"flavor"`
	BackupTime      string     `json:"backupTime"`
	MaintenanceTime string     `json:"maintenanceTime"`
}

func getDatabaseInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(Database)
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_database.getDatabaseInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/cloud/project/%s/database/service/%s", projectId, database.ID), &database)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_database.getDatabaseInfo", err)
		return nil, err
	}
	return database, nil
}

func listDatabase(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_database.listDatabaseInfo", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	var databaseIds []string
	err = client.Get(fmt.Sprintf("/cloud/project/%s/database/service", projectId), &databaseIds)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_database.listDatabaseInfo", err)
		return nil, err
	}
	for _, databaseId := range databaseIds {
		var database Database
		database.ID = databaseId
		d.StreamListItem(ctx, database)
	}
	return nil, nil
}

func getDatabase(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetStringValue()
	var database Database
	database.ID = id
	return database, nil
}
