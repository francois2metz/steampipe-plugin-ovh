package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudAINotebook() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_ai_notebook",
		Description: "OVHcloud AI Notebook gives a quick and simple start launching your Jupyter or VS Code notebooks in the cloud.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listAINotebook,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getAINotebook,
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
				Description: "UUID of the notebook.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Name"),
				Description: "Name of the notebook.",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Region"),
				Description: "Region of the job.",
			},
			{
				Name:        "framework",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Env.FrameworkID"),
				Description: "Framework used by the notebook.",
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Env.FrameworkVersion"),
				Description: "Framework version used by the notebook.",
			},
			{
				Name:        "editor",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Env.EditorID"),
				Description: "Editor used by the notebook.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date when the notebook was created.",
			},
			{
				Name:        "state",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.State"),
				Description: "State of the notebook.",
			},
			{
				Name:        "url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.URL"),
				Description: "Access URL of the notebook.",
			},
		},
	}
}

type AINotebook struct {
	ID        string           `json:"id"`
	CreatedAt time.Time        `json:"createdAt"`
	Spec      AINotebookSpec   `json:"spec"`
	Status    AINotebookStatus `json:"status"`
}

type AINotebookSpec struct {
	Name   string        `json:"name"`
	Region string        `json:"region"`
	Env    AINotebookEnv `json:"env"`
}

type AINotebookEnv struct {
	FrameworkID      string `json:"frameworkId"`
	FrameworkVersion string `json:"frameworkVersion"`
	EditorID         string `json:"editorId"`
}

type AINotebookStatus struct {
	URL   string `json:"url"`
	State string `json:"state"`
}

func getAINotebook(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	notebook := h.Item.(AINotebook)
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_notebook.getAINotebook", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/cloud/project/%s/ai/notebook/%s", projectId, notebook.ID), &notebook)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_notebook.getAINotebook", err)
		return nil, err
	}
	return notebook, nil
}

func listAINotebook(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_notebook.listAINotebook", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	var notebooks []AINotebook
	err = client.Get(fmt.Sprintf("/cloud/project/%s/ai/notebook", projectId), &notebooks)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_notebook.listAINotebook", err)
		return nil, err
	}
	for _, notebook := range notebooks {
		d.StreamListItem(ctx, notebook)
	}
	return nil, nil
}
