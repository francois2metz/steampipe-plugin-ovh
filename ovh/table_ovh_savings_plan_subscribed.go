package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/ovh/go-ovh/ovh"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// SavingsPlan represents an OVH Cloud Savings Plan
type SavingsPlan struct {
	ID              string     `json:"id"`
	DisplayName     string     `json:"displayName"`
	Status          string     `json:"status"`
	Size            int        `json:"size"`
	Flavour         string     `json:"flavour"`
	Duration        string     `json:"duration"`
	AutoRenew       bool       `json:"autoRenew"`
	PeriodEndAction string     `json:"periodEndAction"`
	StartDate       *time.Time `json:"startDate"`
	EndDate         *time.Time `json:"endDate"`
	Region          string     `json:"region"`
	ProductCode     string     `json:"productCode"`
	ServiceName     string     `json:"serviceName"`
	ProjectID       string     `json:"-"` // Set by hydrate function
	ServiceIDNum    int        `json:"-"` // Set by hydrate function
}

func tableOvhSavingsPlanSubscribed() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_savings_plan_subscribed",
		Description: "List OVH Cloud Savings Plans subscribed for each Public Cloud project.",
		List: &plugin.ListConfig{
			Hydrate: listOvhSavingsPlanSubscribed,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "project_id", Require: plugin.Required},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "savings_plan_id"}),
			Hydrate:    getOvhSavingsPlanSubscribed,
		},
		Columns: []*plugin.Column{
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Description: "OVH Public Cloud project ID.",
			},
			{
				Name:        "service_id",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServiceIDNum"),
				Description: "OVH service ID (internal billing ID).",
			},
			{
				Name:        "savings_plan_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
				Description: "Savings plan unique ID.",
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
				Description: "Human-readable plan name.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Plan status (active, terminated, etc.).",
			},
			{
				Name:        "size",
				Type:        proto.ColumnType_INT,
				Description: "Number of resources covered by plan.",
			},
			{
				Name:        "flavour",
				Type:        proto.ColumnType_STRING,
				Description: "Resource type or flavor.",
			},
			{
				Name:        "duration",
				Type:        proto.ColumnType_STRING,
				Description: "Commitment period (ISO-8601, e.g., P12M).",
			},
			{
				Name:        "auto_renewal",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AutoRenew"),
				Description: "Whether plan auto-renews.",
			},
			{
				Name:        "period_end_action",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PeriodEndAction"),
				Description: "What happens at end of commitment (terminate or renew).",
			},
			{
				Name:        "start_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("StartDate"),
				Description: "Plan start date.",
			},
			{
				Name:        "end_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EndDate"),
				Description: "Plan end date.",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Description: "Region (if applicable).",
			},
			{
				Name:        "product_code",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductCode"),
				Description: "Internal plan code (compute_spot_1y, etc.).",
			},
			{
				Name:        "linked_service_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceName"),
				Description: "Friendly name of associated Public Cloud project.",
			},
		},
	}
}

// getServiceIdFromProjectId converts a cloud project UUID to its numeric service ID
func getServiceIdFromProjectId(ctx context.Context, client *ovh.Client, projectId string) (int, error) {
	type serviceInfo struct {
		ServiceId int `json:"serviceId"`
	}

	var service serviceInfo
	err := client.Get(fmt.Sprintf("/cloud/project/%s/serviceInfos", projectId), &service)
	if err != nil {
		return 0, fmt.Errorf("failed to get service ID for project %s: %w", projectId, err)
	}

	return service.ServiceId, nil
}

func listOvhSavingsPlanSubscribed(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.listOvhSavingsPlanSubscribed", "connection_error", err)
		return nil, err
	}

	// Convert project ID to service ID
	serviceID, err := getServiceIdFromProjectId(ctx, client, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.listOvhSavingsPlanSubscribed", "service_id_resolution_error", err)
		return nil, err
	}

	// First get the list of savings plan IDs
	var savingsPlanIDs []string
	err = client.Get(fmt.Sprintf("/services/%d/savingsPlans/subscribed", serviceID), &savingsPlanIDs)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.listOvhSavingsPlanSubscribed", "api_error", err)
		return nil, err
	}

	// Then get the details for each savings plan
	for _, savingsPlanID := range savingsPlanIDs {
		var savingsPlan SavingsPlan
		err = client.Get(fmt.Sprintf("/services/%d/savingsPlans/subscribed/%s", serviceID, savingsPlanID), &savingsPlan)
		if err != nil {
			plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.listOvhSavingsPlanSubscribed", "api_error", err, "savings_plan_id", savingsPlanID)
			continue // Skip this savings plan and continue with others
		}

		// Set the project_id and service_id for the response
		savingsPlan.ProjectID = projectID
		savingsPlan.ServiceIDNum = serviceID

		d.StreamListItem(ctx, savingsPlan)
	}

	return nil, nil
}

func getOvhSavingsPlanSubscribed(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := d.EqualsQuals["project_id"].GetStringValue()
	savingsPlanID := d.EqualsQuals["savings_plan_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.getOvhSavingsPlanSubscribed", "connection_error", err)
		return nil, err
	}

	// Convert project ID to service ID
	serviceID, err := getServiceIdFromProjectId(ctx, client, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.getOvhSavingsPlanSubscribed", "service_id_resolution_error", err)
		return nil, err
	}

	var savingsPlan SavingsPlan
	err = client.Get(fmt.Sprintf("/services/%d/savingsPlans/subscribed/%s", serviceID, savingsPlanID), &savingsPlan)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.getOvhSavingsPlanSubscribed", "api_error", err)
		return nil, err
	}

	// Set the project_id and service_id for the response
	savingsPlan.ProjectID = projectID
	savingsPlan.ServiceIDNum = serviceID

	return savingsPlan, nil
}
