package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudDataJob() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_data_job",
		Description: "A data job is processed by OVH by Apache Spark.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listDataJob,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getDataJob,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{Func: getDataJobInfo},
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
				Hydrate:     getDataJobInfo,
				Type:        proto.ColumnType_STRING,
				Description: "name of the job.",
			},
			{
				Name:        "region",
				Hydrate:     getDataJobInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Openstack region of the job.",
			},
			{
				Name:        "container_name",
				Hydrate:     getDataJobInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Name of the container where the code and the log of the job is.",
			},
			{
				Name:        "engine",
				Hydrate:     getDataJobInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Engine of the job.",
			},
			{
				Name:        "engine_version",
				Hydrate:     getDataJobInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Version of the engine.",
			},
			{
				Name:        "started_at",
				Hydrate:     getDataJobInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Start date of the job.",
				Transform:   transform.FromField("StartDate"),
			},
			{
				Name:        "ended_at",
				Hydrate:     getDataJobInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "End date of the job.",
				Transform:   transform.FromField("EndDate"),
			},
			{
				Name:        "created_at",
				Hydrate:     getDataJobInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Creation date of the job.",
				Transform:   transform.FromField("CreationDate"),
			},
			{
				Name:        "status",
				Hydrate:     getDataJobInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Current state of the job.",
			},
			{
				Name:        "ttl",
				Hydrate:     getDataJobInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Maximum 'Time To Live' (in RFC3339 (duration)) of this job, after which it will be automatically terminated.",
			},
		},
	}
}

type Job struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Region        string    `json:"region"`
	ContainerName string    `json:"containerName"`
	Engine        string    `json:"engine"`
	EngineVersion string    `json:"engineVersion"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	CreationDate  time.Time `json:"creationDate"`
	Status        string    `json:"status"`
	TTL           string    `json:"ttl"`
}

func getDataJobInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	job := h.Item.(Job)
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_data_job.getDataJobInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/cloud/project/%s/dataProcessing/jobs/%s", projectId, job.ID), &job)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_data_job.getDataJobInfo", err)
		return nil, err
	}
	return job, nil
}

func listDataJob(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_data_job.listDataJobInfo", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	var jobIds []string
	err = client.Get(fmt.Sprintf("/cloud/project/%s/dataProcessing/jobs", projectId), &jobIds)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_data_job.listDataJobInfo", err)
		return nil, err
	}
	for _, jobId := range jobIds {
		var job Job
		job.ID = jobId
		d.StreamListItem(ctx, job)
	}
	return nil, nil
}

func getDataJob(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetStringValue()
	var job Job
	job.ID = id
	return job, nil
}
