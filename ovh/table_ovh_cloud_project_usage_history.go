package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type CloudProjectUsageHistory struct {
	ID             string                 `json:"id"`
	ProjectID      string                 `json:"projectId"`
	LastUpdate     *time.Time             `json:"lastUpdate"`
	HourlyUsage    map[string]interface{} `json:"hourlyUsage"`
	MonthlyUsage   map[string]interface{} `json:"monthlyUsage"`
	ResourcesUsage []interface{}          `json:"resourcesUsage"`
	Period         map[string]interface{} `json:"period"`
	UsableCredits  map[string]interface{} `json:"usableCredits"`
}

type UsageHistoryListItem struct {
	ID string `json:"id"`
}

func tableOvhCloudProjectUsageHistory() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_project_usage_history",
		Description: "OVH Cloud Project Usage History with comprehensive resource breakdown and price calculations across multiple billing periods.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listCloudProjectUsageHistory,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "usage_id"}),
			Hydrate:    getCloudProjectUsageHistory,
		},
		Columns: []*plugin.Column{
			// Basic identification
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("project_id"),
				Description: "The project ID.",
			},
			{
				Name:        "usage_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
				Description: "Unique usage period identifier (e.g., RUN2_202511).",
			},
			{
				Name:        "last_update",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("LastUpdate"),
				Description: "Last update timestamp for historical usage data.",
			},

			// Volume storage - historical usage and costs
			{
				Name:        "volumes_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("volume")),
				Description: "Detailed volume usage data (array of historical volume usage with pricing).",
			},
			{
				Name:        "total_volumes_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("volume", OperationTotalPrice)),
				Description: "Total historical price for volume storage.",
			},
			{
				Name:        "volumes_count",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("volume", OperationCount)),
				Description: "Number of different volume types in historical usage.",
			},

			// Instance compute - historical usage and costs
			{
				Name:        "instances_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("instance")),
				Description: "Detailed instance usage data (array of historical instance usage with pricing and savings plans).",
			},
			{
				Name:        "total_instances_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("instance", OperationTotalPrice)),
				Description: "Total historical price for compute instances.",
			},
			{
				Name:        "instances_count",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("instance", OperationCount)),
				Description: "Number of different instance types in historical usage.",
			},

			// Object storage - historical usage and costs
			{
				Name:        "storage_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("storage")),
				Description: "Detailed storage usage data (array of historical object storage usage with pricing).",
			},
			{
				Name:        "total_storage_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("storage", OperationTotalPrice)),
				Description: "Total historical price for object storage.",
			},
			{
				Name:        "storage_count",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("storage", OperationCount)),
				Description: "Number of different storage types in historical usage.",
			},

			// Snapshot storage - historical usage and costs
			{
				Name:        "snapshots_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("snapshot")),
				Description: "Detailed snapshot usage data (array of historical snapshot usage with pricing).",
			},
			{
				Name:        "total_snapshots_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("snapshot", OperationTotalPrice)),
				Description: "Total historical price for volume snapshots.",
			},
			{
				Name:        "snapshots_count",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("snapshot", OperationCount)),
				Description: "Number of different snapshot types in historical usage.",
			},

			// Instance options (floating IPs, additional storage, etc.) - historical usage and costs
			{
				Name:        "instance_options_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("instanceOption")),
				Description: "Detailed instance options usage data (floating IPs, additional resources).",
			},
			{
				Name:        "total_instance_options_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("instanceOption", OperationTotalPrice)),
				Description: "Total historical price for instance options.",
			},

			// Instance bandwidth - historical usage and costs
			{
				Name:        "instance_bandwidth_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("instanceBandwidth")),
				Description: "Detailed instance bandwidth usage data (network usage).",
			},
			{
				Name:        "total_instance_bandwidth_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("instanceBandwidth", OperationTotalPrice)),
				Description: "Total historical price for instance bandwidth.",
			},

			// Managed Kubernetes service - historical usage and costs
			{
				Name:        "kubernetes_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("managedKubernetesService")),
				Description: "Detailed Kubernetes service usage data (managed cluster usage).",
			},
			{
				Name:        "total_kubernetes_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("managedKubernetesService", OperationTotalPrice)),
				Description: "Total historical price for managed Kubernetes services.",
			},

			// Rancher service - historical usage and costs
			{
				Name:        "rancher_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("rancher")),
				Description: "Detailed Rancher service usage data (container management).",
			},
			{
				Name:        "total_rancher_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("rancher", OperationTotalPrice)),
				Description: "Total historical price for Rancher services.",
			},

			// Quantum (AI/ML) - historical usage and costs
			{
				Name:        "quantum_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("quantum")),
				Description: "Detailed quantum/AI usage data (contains nested structure with notebook array).",
			},
			{
				Name:        "total_quantum_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateHistoricalQuantumPrice),
				Description: "Total historical price for quantum/AI services (notebooks, jobs, etc.).",
			},

			// Grand total - overall historical costs
			{
				Name:        "grand_total_historical_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateHistoricalGrandTotal),
				Description: "Grand total of all historical usage costs across all resource types.",
			},

			// Raw hourly usage for complex queries
			{
				Name:        "hourly_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("HourlyUsage"),
				Description: "Complete raw historical usage data structure for advanced analysis.",
			},

			// Monthly usage (savings plans, etc.)
			{
				Name:        "monthly_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("MonthlyUsage"),
				Description: "Monthly usage including savings plan usage.",
			},
			{
				Name:        "monthly_savings_plan_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("MonthlyUsage").Transform(extractUsageField("savingsPlan")),
				Description: "Savings plan monthly usage data with detailed pricing.",
			},
			{
				Name:        "total_monthly_savings_plan_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("MonthlyUsage").Transform(calculateMonthlySavingsPlanPrice),
				Description: "Total historical price for monthly savings plans.",
			},

			// Resources usage (gateways, load balancers, floating IPs, etc.)
			{
				Name:        "resources_usage",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("ResourcesUsage"),
				Description: "Infrastructure resources usage (gateways, load balancers, floating IPs).",
			},
			{
				Name:        "total_resources_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcesUsagePrice),
				Description: "Total historical price for infrastructure resources.",
			},
			{
				Name:        "gateway_historical_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("gateway")),
				Description: "Historical price for gateway resources.",
			},
			{
				Name:        "publicip_historical_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("publicip")),
				Description: "Historical price for public IP resources.",
			},
			{
				Name:        "loadbalancer_historical_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("octavia-loadbalancer")),
				Description: "Historical price for load balancer resources.",
			},
			{
				Name:        "floatingip_historical_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("floatingip")),
				Description: "Historical price for floating IP resources.",
			},

			// Period information
			{
				Name:        "usage_period",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("Period"),
				Description: "Usage period with from/to dates.",
			},
			{
				Name:        "usage_period_from",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("Period.from"),
				Description: "Start date of the usage period.",
			},
			{
				Name:        "usage_period_to",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("Period.to"),
				Description: "End date of the usage period.",
			},

			// Credits information
			{
				Name:        "usable_credits",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("UsableCredits"),
				Description: "Available credits that can be applied to historical usage.",
			},
			{
				Name:        "total_usable_credit",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromField("UsableCredits.totalCredit"),
				Description: "Total amount of usable credits.",
			},

			// Comprehensive grand total
			{
				Name:        "comprehensive_total_historical_price",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getCloudProjectUsageHistory,
				Transform:   transform.FromValue().Transform(calculateComprehensiveHistoricalTotalPrice),
				Description: "Grand total of all historical costs (hourly + monthly + resources).",
			},
		},
	}
}

func listCloudProjectUsageHistory(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	projectID := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCloudProjectUsageHistory.connect", "connection_error", err)
		return nil, err
	}

	// First API call: Get list of usage history records with basic metadata
	endpoint := fmt.Sprintf("/cloud/project/%s/usage/history", projectID)
	var usageRecords []CloudProjectUsageHistory

	if err := client.Get(endpoint, &usageRecords); err != nil {
		plugin.Logger(ctx).Error("listCloudProjectUsageHistory", "api_error", err)
		return nil, err
	}

	// For each usage record, set the project ID and stream the item
	// The detailed data will be hydrated by getCloudProjectUsageHistory
	for _, usage := range usageRecords {
		usage.ProjectID = projectID
		d.StreamListItem(ctx, usage)
	}

	return nil, nil
}

func getCloudProjectUsageHistory(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var projectID, usageID string

	if h.Item != nil {
		// Called from list hydrate
		usage := h.Item.(CloudProjectUsageHistory)
		projectID = usage.ProjectID
		usageID = usage.ID
	} else {
		// Called directly via get
		projectID = d.EqualsQuals["project_id"].GetStringValue()
		usageID = d.EqualsQuals["usage_id"].GetStringValue()
	}

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getCloudProjectUsageHistory.connect", "connection_error", err)
		return nil, err
	}

	// API call: Get detailed usage data
	endpoint := fmt.Sprintf("/cloud/project/%s/usage/history/%s", projectID, usageID)
	var response CloudProjectUsageHistory

	if err := client.Get(endpoint, &response); err != nil {
		plugin.Logger(ctx).Error("getCloudProjectUsageHistory", "usage_id", usageID, "api_error", err)
		return nil, err
	}

	// Set the IDs since the API response includes them in different fields
	response.ProjectID = projectID
	// The API response has "id" field, but we want it mapped to our ID field
	if response.ID == "" {
		response.ID = usageID
	}

	return response, nil
}

// Helper function to calculate quantum price (handles nested structure)
func calculateHistoricalQuantumPrice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return 0.0, nil
	}

	hourlyUsage, ok := d.Value.(map[string]interface{})
	if !ok {
		return 0.0, nil
	}

	quantumInterface, exists := hourlyUsage["quantum"]
	if !exists {
		return 0.0, nil
	}

	quantum, ok := quantumInterface.(map[string]interface{})
	if !ok {
		return 0.0, nil
	}

	// Check notebook array
	if notebookInterface, exists := quantum["notebook"]; exists {
		if notebookList, ok := notebookInterface.([]interface{}); ok {
			return extractTotalPriceFromArray(notebookList, "totalPrice"), nil
		}
	}

	return 0.0, nil
}

// Helper function to calculate grand total across all resource types
func calculateHistoricalGrandTotal(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return 0.0, nil
	}

	hourlyUsage, ok := d.Value.(map[string]interface{})
	if !ok {
		return 0.0, nil
	}

	var grandTotal float64

	// Calculate total for standard resource types
	for _, resourceType := range hourlyResourceTypes {
		if resourceArray, exists := hourlyUsage[resourceType]; exists {
			if resourceList, ok := resourceArray.([]interface{}); ok {
				grandTotal += extractTotalPriceFromArray(resourceList, "totalPrice")
			}
		}
	}

	// Add quantum pricing (special nested structure)
	if quantumInterface, exists := hourlyUsage["quantum"]; exists {
		if quantum, ok := quantumInterface.(map[string]interface{}); ok {
			if notebookInterface, exists := quantum["notebook"]; exists {
				if notebookList, ok := notebookInterface.([]interface{}); ok {
					grandTotal += extractTotalPriceFromArray(notebookList, "totalPrice")
				}
			}
		}
	}

	return grandTotal, nil
}

// Helper function to calculate comprehensive total price (all historical sources)
func calculateComprehensiveHistoricalTotalPrice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return 0.0, nil
	}

	response, ok := d.Value.(CloudProjectUsageHistory)
	if !ok {
		return 0.0, nil
	}

	var grandTotal float64

	// Helper function to safely add calculation results
	addToTotal := func(calculationFunc func(context.Context, *transform.TransformData) (interface{}, error), data interface{}) {
		if result, err := calculationFunc(ctx, &transform.TransformData{Value: data}); err == nil {
			if total, ok := result.(float64); ok {
				grandTotal += total
			}
		}
	}

	// Add all historical totals
	addToTotal(calculateHistoricalGrandTotal, response.HourlyUsage)
	addToTotal(calculateMonthlySavingsPlanPrice, response.MonthlyUsage)
	addToTotal(calculateResourcesUsagePrice, response.ResourcesUsage)

	return grandTotal, nil
}
