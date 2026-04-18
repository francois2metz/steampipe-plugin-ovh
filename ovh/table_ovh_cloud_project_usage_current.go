package ovh

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type CloudProjectUsageCurrent struct {
	ProjectID      string                 `json:"projectId"`
	LastUpdate     string                 `json:"lastUpdate"`
	HourlyUsage    map[string]interface{} `json:"hourlyUsage"`
	MonthlyUsage   map[string]interface{} `json:"monthlyUsage"`
	ResourcesUsage []interface{}          `json:"resourcesUsage"`
	Period         map[string]interface{} `json:"period"`
	UsableCredits  map[string]interface{} `json:"usableCredits"`
}

func tableOvhCloudProjectUsageCurrent() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_project_usage_current",
		Description: "OVH Cloud Project Current Usage with comprehensive resource breakdown and price calculations.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listCloudProjectUsageCurrent,
		},
		Columns: []*plugin.Column{
			// Basic identification
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectID"),
				Description: "The project ID.",
			},
			{
				Name:        "last_update",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastUpdate"),
				Description: "Last update timestamp for current usage data.",
			},

			// Volume storage - current usage and costs
			{
				Name:        "volumes_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("volume")),
				Description: "Detailed volume usage data (array of current volume usage with pricing).",
			},
			{
				Name:        "total_volumes_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("volume", OperationTotalPrice)),
				Description: "Total current price for volume storage.",
			},
			{
				Name:        "volumes_count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("volume", OperationCount)),
				Description: "Number of different volume types in current usage.",
			},

			// Instance compute - current usage and costs
			{
				Name:        "instances_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("instance")),
				Description: "Detailed instance usage data (array of current instance usage with pricing and savings plans).",
			},
			{
				Name:        "total_instances_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("instance", OperationTotalPrice)),
				Description: "Total current price for compute instances.",
			},
			{
				Name:        "instances_count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("instance", OperationCount)),
				Description: "Number of different instance types in current usage.",
			},

			// Object storage - current usage and costs
			{
				Name:        "storage_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("storage")),
				Description: "Detailed storage usage data (array of current object storage usage with pricing).",
			},
			{
				Name:        "total_storage_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("storage", OperationTotalPrice)),
				Description: "Total current price for object storage.",
			},
			{
				Name:        "storage_count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("storage", OperationCount)),
				Description: "Number of different storage types in current usage.",
			},

			// Snapshot storage - current usage and costs
			{
				Name:        "snapshots_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("snapshot")),
				Description: "Detailed snapshot usage data (array of current snapshot usage with pricing).",
			},
			{
				Name:        "total_snapshots_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("snapshot", OperationTotalPrice)),
				Description: "Total current price for volume snapshots.",
			},
			{
				Name:        "snapshots_count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("snapshot", OperationCount)),
				Description: "Number of different snapshot types in current usage.",
			},

			// Instance options (floating IPs, additional storage, etc.) - current usage and costs
			{
				Name:        "instance_options_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("instanceOption")),
				Description: "Detailed instance options usage data (floating IPs, additional resources).",
			},
			{
				Name:        "total_instance_options_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("instanceOption", OperationTotalPrice)),
				Description: "Total current price for instance options.",
			},

			// Instance bandwidth - current usage and costs
			{
				Name:        "instance_bandwidth_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("instanceBandwidth")),
				Description: "Detailed instance bandwidth usage data (network usage).",
			},
			{
				Name:        "total_instance_bandwidth_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("instanceBandwidth", OperationTotalPrice)),
				Description: "Total current price for instance bandwidth.",
			},

			// Managed Kubernetes service - current usage and costs
			{
				Name:        "kubernetes_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("managedKubernetesService")),
				Description: "Detailed Kubernetes service usage data (managed cluster usage).",
			},
			{
				Name:        "total_kubernetes_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("managedKubernetesService", OperationTotalPrice)),
				Description: "Total current price for managed Kubernetes services.",
			},

			// Rancher service - current usage and costs
			{
				Name:        "rancher_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("rancher")),
				Description: "Detailed Rancher service usage data (container management).",
			},
			{
				Name:        "total_rancher_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(extractHourlyUsageField("rancher", OperationTotalPrice)),
				Description: "Total current price for Rancher services.",
			},

			// Quantum (AI/ML) - current usage and costs
			{
				Name:        "quantum_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("quantum")),
				Description: "Detailed quantum/AI usage data (contains nested structure with notebook array).",
			},
			{
				Name:        "total_quantum_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateCurrentQuantumPrice),
				Description: "Total current price for quantum/AI services (notebooks, jobs, etc.).",
			},

			// Grand total - overall current costs
			{
				Name:        "grand_total_current_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateCurrentGrandTotal),
				Description: "Grand total of all current usage costs across all resource types.",
			},

			// Raw hourly usage for complex queries
			{
				Name:        "hourly_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage"),
				Description: "Complete raw current usage data structure for advanced analysis.",
			},

			// Monthly usage (savings plans, etc.)
			{
				Name:        "monthly_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MonthlyUsage"),
				Description: "Monthly usage including savings plan usage.",
			},
			{
				Name:        "monthly_savings_plan_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MonthlyUsage").Transform(extractUsageField("savingsPlan")),
				Description: "Savings plan monthly usage data with detailed pricing.",
			},
			{
				Name:        "total_monthly_savings_plan_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("MonthlyUsage").Transform(calculateMonthlySavingsPlanPrice),
				Description: "Total current price for monthly savings plans.",
			},

			// Resources usage (gateways, load balancers, floating IPs, etc.)
			{
				Name:        "resources_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourcesUsage"),
				Description: "Infrastructure resources usage (gateways, load balancers, floating IPs).",
			},
			{
				Name:        "total_resources_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcesUsagePrice),
				Description: "Total current price for infrastructure resources.",
			},
			{
				Name:        "gateway_current_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("gateway")),
				Description: "Current price for gateway resources.",
			},
			{
				Name:        "publicip_current_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("publicip")),
				Description: "Current price for public IP resources.",
			},
			{
				Name:        "loadbalancer_current_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("octavia-loadbalancer")),
				Description: "Current price for load balancer resources.",
			},
			{
				Name:        "floatingip_current_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("floatingip")),
				Description: "Current price for floating IP resources.",
			},

			// Period information
			{
				Name:        "usage_period",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Period"),
				Description: "Usage period with from/to dates.",
			},
			{
				Name:        "usage_period_from",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Period.from"),
				Description: "Start date of the usage period.",
			},
			{
				Name:        "usage_period_to",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Period.to"),
				Description: "End date of the usage period.",
			},

			// Credits information
			{
				Name:        "usable_credits",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("UsableCredits"),
				Description: "Available credits that can be applied to current usage.",
			},
			{
				Name:        "total_usable_credit",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("UsableCredits.totalCredit"),
				Description: "Total amount of usable credits.",
			},

			// Comprehensive grand total
			{
				Name:        "comprehensive_total_current_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromValue().Transform(calculateComprehensiveCurrentTotalPrice),
				Description: "Grand total of all current costs (hourly + monthly + resources).",
			},
		},
	}
}

func listCloudProjectUsageCurrent(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	projectID := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCloudProjectUsageCurrent.connect", "connection_error", err)
		return nil, err
	}

	endpoint := fmt.Sprintf("/cloud/project/%s/usage/current", projectID)
	var response CloudProjectUsageCurrent

	if err := client.Get(endpoint, &response); err != nil {
		plugin.Logger(ctx).Error("listCloudProjectUsageCurrent", "api_error", err)
		return nil, err
	}

	// Set the project ID since the API response doesn't include it
	response.ProjectID = projectID

	d.StreamListItem(ctx, response)
	return nil, nil
}

// Helper function to extract usage field (reusing from forecast table pattern)
func extractUsageField(fieldName string) func(context.Context, *transform.TransformData) (interface{}, error) {
	return func(ctx context.Context, d *transform.TransformData) (interface{}, error) {
		if d.Value == nil {
			return nil, nil
		}

		hourlyUsage, ok := d.Value.(map[string]interface{})
		if !ok {
			return nil, nil
		}

		field, exists := hourlyUsage[fieldName]
		if !exists {
			return nil, nil
		}

		return field, nil
	}
}

// Helper function to calculate quantum price (handles nested structure)
func calculateCurrentQuantumPrice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func calculateCurrentGrandTotal(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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

// Helper function to calculate comprehensive total price (all current sources)
func calculateComprehensiveCurrentTotalPrice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return 0.0, nil
	}

	response, ok := d.Value.(CloudProjectUsageCurrent)
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

	// Add all current totals
	addToTotal(calculateCurrentGrandTotal, response.HourlyUsage)
	addToTotal(calculateMonthlySavingsPlanPrice, response.MonthlyUsage)
	addToTotal(calculateResourcesUsagePrice, response.ResourcesUsage)

	return grandTotal, nil
}
