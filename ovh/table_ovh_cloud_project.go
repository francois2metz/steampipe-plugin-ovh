package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableOvhCloudProject() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_project",
		Description: "List available project.",
		List: &plugin.ListConfig{
			Hydrate: listProject,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getProject,
		},
		HydrateDependencies: []plugin.HydrateDependencies{
			{Func: getProjectInfo},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Project ID.",
			},
			{
				Name:        "name",
				Hydrate:     getProjectInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Project name.",
			},
			{
				Name:        "description",
				Hydrate:     getProjectInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Project description.",
			},
			{
				Name:        "plan_code",
				Hydrate:     getProjectInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Order plan code.",
			},
			{
				Name:        "order_id",
				Hydrate:     getProjectInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Project order ID.",
			},
			{
				Name:        "status",
				Hydrate:     getProjectInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Project status (creating, deleted, deleting, ok, suspended)",
			},
			{
				Name:        "unleash",
				Hydrate:     getProjectInfo,
				Type:        proto.ColumnType_BOOL,
				Description: "Project unleashed.",
			},
			{
				Name:        "manual_quota",
				Hydrate:     getProjectInfo,
				Type:        proto.ColumnType_BOOL,
				Description: "Manual quota prevent automatic quota upgrade.",
			},
			{
				Name:        "expiration_at",
				Hydrate:     getProjectInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Expiration"),
				Description: "Expiration date of your project. After this date, your project will be deleted.",
			},
			{
				Name:        "created_at",
				Hydrate:     getProjectInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationDate"),
				Description: "Project creation date.",
			},
		},
	}
}

type Project struct {
	ID           string     `json:"project_id"`
	Name         string     `json:"projectName"`
	Description  string     `json:"description"`
	PlanCode     string     `json:"planCode"`
	Unleash      *bool      `json:"unleash"`
	Expiration   *time.Time `json:"expiration,omitempty"`
	CreationDate time.Time  `json:"creationDate"`
	OrderId      int        `json:"orderId"`
	Access       string     `json:"access"`
	Status       string     `json:"status"`
	ManualQuota  *bool      `json:"manualQuota"`
}

func getProjectInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getProjectInfo")
	project := h.Item.(Project)

	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/cloud/project/%s", project.ID), &project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func listProject(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	var projects []string
	err = client.Get("/cloud/project", &projects)
	if err != nil {
		return nil, err
	}
	for _, projectId := range projects {
		var project Project
		project.ID = projectId
		d.StreamListItem(ctx, project)
	}
	return nil, nil
}

func getProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	projectId := quals["id"].GetStringValue()
	var project Project
	project.ID = projectId
	return project, nil
}
