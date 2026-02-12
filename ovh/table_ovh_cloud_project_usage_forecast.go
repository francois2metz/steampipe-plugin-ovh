package ovh

import (
	"context"
	"fmt"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Constants for forecast operations
const (
	OperationTotalPrice = "totalPrice"
	OperationCount      = "count"
)

// Resource type mapping for grand total calculation
var hourlyResourceTypes = []string{
	"volume", "instance", "storage", "snapshot",
	"instanceOption", "instanceBandwidth",
	"managedKubernetesService", "rancher",
}

type CloudProjectUsageForecast struct {
	ProjectID      string                 `json:"projectId"`
	LastUpdate     string                 `json:"lastUpdate"`
	HourlyUsage    map[string]interface{} `json:"hourlyUsage"`
	MonthlyUsage   map[string]interface{} `json:"monthlyUsage"`
	ResourcesUsage []interface{}          `json:"resourcesUsage"`
	Period         map[string]interface{} `json:"period"`
	UsableCredits  map[string]interface{} `json:"usableCredits"`
}

func tableOvhCloudProjectUsageForecast() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_project_usage_forecast",
		Description: "OVH Cloud Project Usage Forecast with comprehensive resource breakdown and price calculations.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listCloudProjectUsageForecast,
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
				Description: "Last update timestamp for forecast data.",
			},

			// Volume storage - predicted usage and costs
			{
				Name:        "volumes_forecast",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("volume")),
				Description: "Detailed volume forecast data (array of predicted volume usage with pricing).",
			},
			{
				Name:        "total_volumes_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateForecastTotalPrice("volume")),
				Description: "Total forecasted price for volume storage.",
			},
			{
				Name:        "volumes_count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateUsageCount("volume")),
				Description: "Number of different volume types in forecast.",
			},

			// Instance compute - predicted usage and costs
			{
				Name:        "instances_forecast",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("instance")),
				Description: "Detailed instance forecast data (array of predicted instance usage with pricing and savings plans).",
			},
			{
				Name:        "total_instances_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateForecastTotalPrice("instance")),
				Description: "Total forecasted price for compute instances.",
			},
			{
				Name:        "instances_count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateUsageCount("instance")),
				Description: "Number of different instance types in forecast.",
			},

			// Object storage - predicted usage and costs
			{
				Name:        "storage_forecast",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("storage")),
				Description: "Detailed storage forecast data (array of predicted object storage usage with pricing).",
			},
			{
				Name:        "total_storage_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateForecastTotalPrice("storage")),
				Description: "Total forecasted price for object storage.",
			},
			{
				Name:        "storage_count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateUsageCount("storage")),
				Description: "Number of different storage types in forecast.",
			},

			// Snapshot storage - predicted usage and costs
			{
				Name:        "snapshots_forecast",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("snapshot")),
				Description: "Detailed snapshot forecast data (array of predicted snapshot usage with pricing).",
			},
			{
				Name:        "total_snapshots_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateForecastTotalPrice("snapshot")),
				Description: "Total forecasted price for volume snapshots.",
			},
			{
				Name:        "snapshots_count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateUsageCount("snapshot")),
				Description: "Number of different snapshot types in forecast.",
			},

			// Instance options (floating IPs, additional storage, etc.) - predicted usage and costs
			{
				Name:        "instance_options_forecast",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("instanceOption")),
				Description: "Detailed instance options forecast data (floating IPs, additional resources).",
			},
			{
				Name:        "total_instance_options_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateForecastTotalPrice("instanceOption")),
				Description: "Total forecasted price for instance options.",
			},

			// Instance bandwidth - predicted usage and costs
			{
				Name:        "instance_bandwidth_forecast",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("instanceBandwidth")),
				Description: "Detailed instance bandwidth forecast data (network usage predictions).",
			},
			{
				Name:        "total_instance_bandwidth_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateForecastTotalPrice("instanceBandwidth")),
				Description: "Total forecasted price for instance bandwidth.",
			},

			// Managed Kubernetes service - predicted usage and costs
			{
				Name:        "kubernetes_forecast",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("managedKubernetesService")),
				Description: "Detailed Kubernetes service forecast data (managed cluster predictions).",
			},
			{
				Name:        "total_kubernetes_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateForecastTotalPrice("managedKubernetesService")),
				Description: "Total forecasted price for managed Kubernetes services.",
			},

			// Rancher service - predicted usage and costs
			{
				Name:        "rancher_forecast",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("rancher")),
				Description: "Detailed Rancher service forecast data (container management predictions).",
			},
			{
				Name:        "total_rancher_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateForecastTotalPrice("rancher")),
				Description: "Total forecasted price for Rancher services.",
			},

			// Quantum (AI/ML) - predicted usage and costs
			{
				Name:        "quantum_forecast",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage").Transform(extractUsageField("quantum")),
				Description: "Detailed quantum/AI forecast data (contains nested structure with notebook array predictions).",
			},
			{
				Name:        "total_quantum_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateForecastQuantumPrice),
				Description: "Total forecasted price for quantum/AI services (notebooks, jobs, etc.).",
			},

			// Grand total - overall predicted costs
			{
				Name:        "grand_total_forecast_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("HourlyUsage").Transform(calculateForecastGrandTotal),
				Description: "Grand total of all forecasted usage costs across all resource types.",
			},

			// Raw hourly usage for complex queries
			{
				Name:        "hourly_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HourlyUsage"),
				Description: "Complete raw forecast data structure for advanced analysis.",
			},

			// Monthly usage forecasts (savings plans, etc.)
			{
				Name:        "monthly_usage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MonthlyUsage"),
				Description: "Monthly usage forecasts including savings plan predictions.",
			},
			{
				Name:        "monthly_savings_plan_forecast",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MonthlyUsage").Transform(extractUsageField("savingsPlan")),
				Description: "Savings plan monthly forecast data with detailed pricing.",
			},
			{
				Name:        "total_monthly_savings_plan_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("MonthlyUsage").Transform(calculateMonthlySavingsPlanPrice),
				Description: "Total forecasted price for monthly savings plans.",
			},

			// Resources usage forecasts (gateways, load balancers, floating IPs, etc.)
			{
				Name:        "resources_usage_forecast",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourcesUsage"),
				Description: "Infrastructure resources usage forecasts (gateways, load balancers, floating IPs).",
			},
			{
				Name:        "total_resources_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcesUsagePrice),
				Description: "Total forecasted price for infrastructure resources.",
			},
			{
				Name:        "gateway_forecast_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("gateway")),
				Description: "Forecasted price for gateway resources.",
			},
			{
				Name:        "publicip_forecast_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("publicip")),
				Description: "Forecasted price for public IP resources.",
			},
			{
				Name:        "loadbalancer_forecast_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("octavia-loadbalancer")),
				Description: "Forecasted price for load balancer resources.",
			},
			{
				Name:        "floatingip_forecast_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ResourcesUsage").Transform(calculateResourcePrice("floatingip")),
				Description: "Forecasted price for floating IP resources.",
			},

			// Forecast period information
			{
				Name:        "forecast_period",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Period"),
				Description: "Forecast period with from/to dates.",
			},
			{
				Name:        "forecast_period_from",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Period.from"),
				Description: "Start date of the forecast period.",
			},
			{
				Name:        "forecast_period_to",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Period.to"),
				Description: "End date of the forecast period.",
			},

			// Credits information
			{
				Name:        "usable_credits",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("UsableCredits"),
				Description: "Available credits that can be applied to forecasted usage.",
			},
			{
				Name:        "total_usable_credit",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("UsableCredits.totalCredit"),
				Description: "Total amount of usable credits.",
			},

			// Comprehensive grand total
			{
				Name:        "comprehensive_total_forecast_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromValue().Transform(calculateComprehensiveTotalPrice),
				Description: "Grand total of all forecasted costs (hourly + monthly + resources).",
			},
		},
	}
}

func listCloudProjectUsageForecast(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	projectID := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCloudProjectUsageForecast.connect", "connection_error", err)
		return nil, err
	}

	endpoint := fmt.Sprintf("/cloud/project/%s/usage/forecast", projectID)
	var response CloudProjectUsageForecast

	if err := client.Get(endpoint, &response); err != nil {
		plugin.Logger(ctx).Error("listCloudProjectUsageForecast", "api_error", err)
		return nil, err
	}

	// Set the project ID since the API response doesn't include it
	response.ProjectID = projectID

	d.StreamListItem(ctx, response)
	return nil, nil
}

// Generic helper function to extract field from hourly usage by resource type
func extractHourlyUsageField(resourceType, operation string) func(context.Context, *transform.TransformData) (interface{}, error) {
	return func(ctx context.Context, d *transform.TransformData) (interface{}, error) {
		if d.Value == nil {
			switch operation {
			case OperationCount:
				return 0, nil
			default:
				return 0.0, nil
			}
		}

		hourlyUsage, ok := d.Value.(map[string]interface{})
		if !ok {
			switch operation {
			case OperationCount:
				return 0, nil
			default:
				return 0.0, nil
			}
		}

		resourceArray, exists := hourlyUsage[resourceType]
		if !exists {
			switch operation {
			case OperationCount:
				return 0, nil
			default:
				return 0.0, nil
			}
		}

		resourceList, ok := resourceArray.([]interface{})
		if !ok {
			switch operation {
			case OperationCount:
				return 0, nil
			default:
				return 0.0, nil
			}
		}

		switch operation {
		case OperationTotalPrice:
			return extractTotalPriceFromArray(resourceList, "totalPrice"), nil
		case OperationCount:
			return len(resourceList), nil
		default:
			return 0.0, nil
		}
	}
}

// Helper function to calculate total price for a resource type
func calculateForecastTotalPrice(resourceType string) func(context.Context, *transform.TransformData) (interface{}, error) {
	return extractHourlyUsageField(resourceType, OperationTotalPrice)
}

// Helper function to calculate count of different resource types
func calculateUsageCount(resourceType string) func(context.Context, *transform.TransformData) (interface{}, error) {
	return extractHourlyUsageField(resourceType, OperationCount)
}

// Helper function to calculate quantum price (handles nested structure)
func calculateForecastQuantumPrice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
func calculateForecastGrandTotal(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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

// Helper function to calculate monthly savings plan price
func calculateMonthlySavingsPlanPrice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return 0.0, nil
	}

	monthlyUsage, ok := d.Value.(map[string]interface{})
	if !ok {
		return 0.0, nil
	}

	savingsPlanInterface, exists := monthlyUsage["savingsPlan"]
	if !exists {
		return 0.0, nil
	}

	savingsPlanList, ok := savingsPlanInterface.([]interface{})
	if !ok {
		return 0.0, nil
	}

	return extractNestedPriceFromArray(savingsPlanList, "totalPrice", "value"), nil
}

// Helper function to calculate total resources usage price
func calculateResourcesUsagePrice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return 0.0, nil
	}

	resourcesUsage, ok := d.Value.([]interface{})
	if !ok {
		return 0.0, nil
	}

	return extractTotalPriceFromArray(resourcesUsage, "totalPrice"), nil
}

// Helper function to calculate price for a specific resource type
func calculateResourcePrice(resourceType string) func(context.Context, *transform.TransformData) (interface{}, error) {
	return func(ctx context.Context, d *transform.TransformData) (interface{}, error) {
		if d.Value == nil {
			return 0.0, nil
		}

		resourcesUsage, ok := d.Value.([]interface{})
		if !ok {
			return 0.0, nil
		}

		return extractPriceFromFilteredArray(resourcesUsage, "type", resourceType, "totalPrice"), nil
	}
}

// Helper function to calculate comprehensive total price (all forecast sources)
func calculateComprehensiveTotalPrice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return 0.0, nil
	}

	response, ok := d.Value.(CloudProjectUsageForecast)
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

	// Add all forecast totals
	addToTotal(calculateForecastGrandTotal, response.HourlyUsage)
	addToTotal(calculateMonthlySavingsPlanPrice, response.MonthlyUsage)
	addToTotal(calculateResourcesUsagePrice, response.ResourcesUsage)

	return grandTotal, nil
}

// Generic helper function to extract and sum prices from an array
func extractTotalPriceFromArray(items []interface{}, priceField string) float64 {
	var totalPrice float64
	for _, item := range items {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if priceInterface, exists := itemMap[priceField]; exists {
				totalPrice += convertToFloat64(priceInterface)
			}
		}
	}
	return totalPrice
}

// Generic helper function to extract and sum nested prices from an array
func extractNestedPriceFromArray(items []interface{}, parentField, childField string) float64 {
	var totalPrice float64
	for _, item := range items {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if parentInterface, exists := itemMap[parentField]; exists {
				if parentMap, ok := parentInterface.(map[string]interface{}); ok {
					if childInterface, exists := parentMap[childField]; exists {
						totalPrice += convertToFloat64(childInterface)
					}
				}
			}
		}
	}
	return totalPrice
}

// Generic helper function to extract prices from filtered array items
func extractPriceFromFilteredArray(items []interface{}, filterField, filterValue, priceField string) float64 {
	var totalPrice float64
	for _, item := range items {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if typeInterface, exists := itemMap[filterField]; exists {
				if typeStr, ok := typeInterface.(string); ok && typeStr == filterValue {
					if priceInterface, exists := itemMap[priceField]; exists {
						totalPrice += convertToFloat64(priceInterface)
					}
				}
			}
		}
	}
	return totalPrice
}

// Generic helper function to convert various numeric types to float64
func convertToFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		// Handle json.Number case
		if strVal := fmt.Sprintf("%v", v); strVal != "" {
			if price, err := strconv.ParseFloat(strVal, 64); err == nil {
				return price
			}
		}
	}
	return 0.0
}
