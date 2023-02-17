package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudProject() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_project",
		Description: "A cloud project is a way to regroup instance, storage, database, ... under a name.",
		List: &plugin.ListConfig{
			Hydrate: listProject,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getProject,
		},
		HydrateConfig: []plugin.HydrateConfig{
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
	project := h.Item.(Project)

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_project.getProjectInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/cloud/project/%s", project.ID), &project)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_project.getProjectInfo", err)
		return nil, err
	}
	return project, nil
}

func listProject(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_project.listProject", "connection_error", err)
		return nil, err
	}
	var projects []string
	err = client.Get("/cloud/project", &projects)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_project.listProject", err)
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
	quals := d.EqualsQuals
	projectId := quals["id"].GetStringValue()
	var project Project
	project.ID = projectId
	return project, nil
}
