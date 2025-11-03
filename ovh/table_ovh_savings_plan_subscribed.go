package ovh

import (
	"context"
	"fmt"
	"strings"

	"github.com/ovh/go-ovh/ovh"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// PlannedChange represents a change planned on a Savings Plan
type PlannedChange struct {
	PlannedOn  string                 `json:"plannedOn"`
	Properties map[string]interface{} `json:"properties"`
}

// SavingsPlan represents an OVH Cloud Savings Plan
type SavingsPlan struct {
	ID              string          `json:"id"`
	DisplayName     string          `json:"displayName"`
	Status          string          `json:"status"`
	Size            int             `json:"size"`
	Flavor          *string         `json:"flavor"`          // Optional field
	Period          *string         `json:"period"`          // Optional field
	OfferID         *string         `json:"offerId"`         // Optional field
	PeriodEndAction *string         `json:"periodEndAction"` // Optional field
	StartDate       *string         `json:"startDate"`       // Optional field
	EndDate         *string         `json:"endDate"`         // Optional field
	PeriodStartDate *string         `json:"periodStartDate"` // Optional field
	PeriodEndDate   *string         `json:"periodEndDate"`   // Optional field
	TerminationDate *string         `json:"terminationDate"` // Optional field - can be null
	PlannedChanges  []PlannedChange `json:"plannedChanges"`  // Array of planned changes
	ProjectID       string          `json:"-"`               // Set by hydrate function
	ServiceIDNum    int             `json:"-"`               // Set by hydrate function
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
				Name:        "flavor",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Flavor"),
				Description: "Savings Plan flavor (resource type).",
			},
			{
				Name:        "period",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Period"),
				Description: "Periodicity of the Savings Plan (duration, e.g., P1Y).",
			},
			{
				Name:        "offer_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OfferID"),
				Description: "Savings Plan commercial offer identifier.",
			},
			{
				Name:        "period_end_action",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PeriodEndAction"),
				Description: "Action performed when reaching the end of the period (REACTIVATE or TERMINATE).",
			},
			{
				Name:        "start_date",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StartDate"),
				Description: "Start date of the Savings Plan.",
			},
			{
				Name:        "end_date",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndDate"),
				Description: "End date of the Savings Plan.",
			},
			{
				Name:        "period_start_date",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PeriodStartDate"),
				Description: "Start date of the current period.",
			},
			{
				Name:        "period_end_date",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PeriodEndDate"),
				Description: "End date of the current period.",
			},
			{
				Name:        "termination_date",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TerminationDate"),
				Description: "Date at which the Savings Plan is scheduled to be terminated (null if not scheduled for termination).",
			},
			{
				Name:        "planned_changes",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PlannedChanges"),
				Description: "Changes planned on the Savings Plan.",
			},
		},
	}
}

// getServiceIDFromProjectID converts a cloud project UUID to its numeric service ID
func getServiceIDFromProjectID(ctx context.Context, client *ovh.Client, projectID string) (int, error) {
	type serviceInfo struct {
		ServiceID *int `json:"serviceId"`
	}

	var service serviceInfo
	err := client.Get(fmt.Sprintf("/cloud/project/%s/serviceInfos", projectID), &service)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.getServiceIDFromProjectID", err)
		return 0, err
	}

	if service.ServiceID == nil {
		return 0, fmt.Errorf("serviceId not found in response for project %s", projectID)
	}
	return *service.ServiceID, nil
}

func listOvhSavingsPlanSubscribed(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.listOvhSavingsPlanSubscribed", "connection_error", err)
		return nil, err
	}

	// Convert project ID to service ID
	serviceID, err := getServiceIDFromProjectID(ctx, client, projectID)
	if err != nil {
		return nil, nil // Return empty result for projects without savings plan support
	}

	// OVH API can return either:
	// 1. Empty array [] when no savings plans exist
	// 2. Array of full objects [{...}] when savings plans exist

	// First try to get as full objects (most common case when savings plans exist)
	var savingsPlans []SavingsPlan
	err = client.Get(fmt.Sprintf("/services/%d/savingsPlans/subscribed", serviceID), &savingsPlans)
	if err != nil {
		// If we get a JSON unmarshal error, it might be that the API returned string IDs instead of objects
		if strings.Contains(err.Error(), "cannot unmarshal") {
			// Try to get as string array (fallback)
			var savingsPlanIDs []string
			err2 := client.Get(fmt.Sprintf("/services/%d/savingsPlans/subscribed", serviceID), &savingsPlanIDs)
			if err2 != nil {
				plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.listOvhSavingsPlanSubscribed", err2)
				return nil, err2
			}

			// Get details for each savings plan ID
			for _, savingsPlanID := range savingsPlanIDs {
				var savingsPlan SavingsPlan
				err3 := client.Get(fmt.Sprintf("/services/%d/savingsPlans/subscribed/%s", serviceID, savingsPlanID), &savingsPlan)
				if err3 != nil {
					plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.listOvhSavingsPlanSubscribed", err3)
					continue // Skip this savings plan and continue with others
				}

				// Set the project_id and service_id for the response
				savingsPlan.ProjectID = projectID
				savingsPlan.ServiceIDNum = serviceID

				d.StreamListItem(ctx, savingsPlan)
			}
			return nil, nil
		}

		// If the API returns 404 or similar, it means no savings plans exist or the service doesn't support them
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			return nil, nil // Return empty result instead of error
		}
		plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.listOvhSavingsPlanSubscribed", err)
		return nil, err
	}

	// Process the savings plans we got as full objects
	for _, savingsPlan := range savingsPlans {
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
	serviceID, err := getServiceIDFromProjectID(ctx, client, projectID)
	if err != nil {
		return nil, err
	}

	var savingsPlan SavingsPlan
	err = client.Get(fmt.Sprintf("/services/%d/savingsPlans/subscribed/%s", serviceID, savingsPlanID), &savingsPlan)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("ovh_savings_plan_subscribed.getOvhSavingsPlanSubscribed", err)
		return nil, err
	}

	// Set the project_id and service_id for the response
	savingsPlan.ProjectID = projectID
	savingsPlan.ServiceIDNum = serviceID

	return savingsPlan, nil
}
