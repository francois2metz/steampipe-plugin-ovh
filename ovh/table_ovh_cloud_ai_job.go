package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudAIJob() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_ai_job",
		Description: "OVHcloud AI Training lets you train your AI, machine learning and deep learning models efficiently and easily, and optimise your GPU usage.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listAIJob,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getAIJob,
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
				Description: "UUID of the job.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Name"),
				Description: "Name of the job.",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Region"),
				Description: "Region of the job.",
			},
			{
				Name:        "image",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Image"),
				Description: "Docker image used by the job.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date when the job was created.",
			},
			{
				Name:        "state",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.State"),
				Description: "State of the job.",
			},
			{
				Name:        "url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.URL"),
				Description: "Access URL of the job.",
			},
		},
	}
}

type AIJob struct {
	ID        string      `json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	Spec      AIJobSpec   `json:"spec"`
	Status    AIJobStatus `json:"status"`
}

type AIJobSpec struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Region string `json:"region"`
}

type AIJobStatus struct {
	URL   string `json:"url"`
	State string `json:"state"`
}

func getAIJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	job := h.Item.(AIJob)
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_job.getAIJob", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/cloud/project/%s/ai/job/%s", projectId, job.ID), &job)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_job.getAIJob", err)
		return nil, err
	}
	return job, nil
}

func listAIJob(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_job.listAIJob", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	var jobs []AIJob
	err = client.Get(fmt.Sprintf("/cloud/project/%s/ai/job", projectId), &jobs)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ai_job.listAIJob", err)
		return nil, err
	}
	for _, job := range jobs {
		d.StreamListItem(ctx, job)
	}
	return nil, nil
}
