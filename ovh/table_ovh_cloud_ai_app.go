package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableOvhCloudAIApp() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_ai_app",
		Description: "OVHcloud AI Deploy lets you easily deploy machine learning models and applications to production, create your API access points effortlessly, and make effective predictions.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listAIApp,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getAIApp,
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
				Description: "UUID of the app.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Name"),
				Description: "Name of the app.",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Region"),
				Description: "Region of the app.",
			},
			{
				Name:        "image",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Image"),
				Description: "Docker image used by the app.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date when the app was created.",
			},
			{
				Name:        "state",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.State"),
				Description: "State of the app.",
			},
			{
				Name:        "replicas",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Status.AvailableReplicas"),
				Description: "Available replicas of the app.",
			},
			{
				Name:        "url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.URL"),
				Description: "Access URL of the app.",
			},
		},
	}
}

type AIApp struct {
	ID        string      `json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	Spec      AIAppSpec   `json:"spec"`
	Status    AIAppStatus `json:"status"`
}

type AIAppSpec struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Region string `json:"region"`
}

type AIAppStatus struct {
	URL               string `json:"url"`
	State             string `json:"state"`
	AvailableReplicas int    `json:"availableReplicas"`
}

func getAIApp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	app := h.Item.(AIApp)
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_app.getAIApp", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/cloud/project/%s/ai/app/%s", projectId, app.ID), &app)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_app.getAIApp", err)
		return nil, err
	}
	return app, nil
}

func listAIApp(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_app.listAIApp", "connection_error", err)
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	var apps []AIApp
	err = client.Get(fmt.Sprintf("/cloud/project/%s/ai/app", projectId), &apps)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_app.listAIApp", err)
		return nil, err
	}
	for _, app := range apps {
		d.StreamListItem(ctx, app)
	}
	return nil, nil
}
