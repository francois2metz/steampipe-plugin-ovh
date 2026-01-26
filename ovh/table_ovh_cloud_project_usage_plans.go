package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOvhCloudProjectUsagePlans() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_project_usage_plans",
		Description: "Savings plan usage and cost information for OVH Cloud projects.",
		List: &plugin.ListConfig{
			Hydrate: listCloudProjectUsagePlans,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "project_id", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectID"),
				Description: "OVH Public Cloud project ID.",
			},
			{
				Name:        "period_from",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Period.From"),
				Description: "Start of the usage period.",
			},
			{
				Name:        "period_to",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Period.To"),
				Description: "End of the usage period.",
			},
			{
				Name:        "total_savings",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("TotalSavings.Value"),
				Description: "Total amount saved by using savings plans.",
			},
			{
				Name:        "total_savings_currency",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TotalSavings.CurrencyCode"),
				Description: "Currency code for total savings.",
			},
			{
				Name:        "total_savings_text",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TotalSavings.Text"),
				Description: "Human-readable total savings amount.",
			},
			// Flavor-level columns (from first flavor in array)
			{
				Name:        "flavor",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Flavors").Transform(extractFirstFlavorName),
				Description: "Instance flavor type (e.g., b3-16).",
			},
			{
				Name:        "flat_fee_total_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Flavors").Transform(extractFlatFeeTotalPrice),
				Description: "Total flat fee price for the flavor.",
			},
			{
				Name:        "flat_fee_currency",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Flavors").Transform(extractFlatFeeCurrency),
				Description: "Currency for flat fee pricing.",
			},
			{
				Name:        "over_quota_quantity",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Flavors").Transform(extractOverQuotaQuantity),
				Description: "Quantity of over-quota usage.",
			},
			{
				Name:        "over_quota_unit_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Flavors").Transform(extractOverQuotaUnitPrice),
				Description: "Unit price for over-quota usage.",
			},
			{
				Name:        "flavor_total_price",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Flavors").Transform(extractFlavorTotalPrice),
				Description: "Total price for this flavor.",
			},
			{
				Name:        "flavor_saved_amount",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Flavors").Transform(extractFlavorSavedAmount),
				Description: "Amount saved for this flavor.",
			},
			// Usage period information
			{
				Name:        "usage_period_coverage",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Flavors").Transform(extractUsageCoverage),
				Description: "Coverage percentage for the usage period.",
			},
			{
				Name:        "usage_period_utilization",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Flavors").Transform(extractUsageUtilization),
				Description: "Utilization percentage for the usage period.",
			},
			{
				Name:        "consumption_size",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Flavors").Transform(extractConsumptionSize),
				Description: "Number of instances consumed.",
			},
			{
				Name:        "cumul_plan_size",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Flavors").Transform(extractCumulPlanSize),
				Description: "Cumulative plan size.",
			},
			// Subscription information
			{
				Name:        "subscription_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Flavors").Transform(extractSubscriptionID),
				Description: "ID of the savings plan subscription.",
			},
			{
				Name:        "subscription_size",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Flavors").Transform(extractSubscriptionSize),
				Description: "Size of the savings plan subscription.",
			},
			{
				Name:        "subscription_begin",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Flavors").Transform(extractSubscriptionBegin),
				Description: "Start date of the savings plan subscription.",
			},
			{
				Name:        "subscription_end",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Flavors").Transform(extractSubscriptionEnd),
				Description: "End date of the savings plan subscription.",
			},
			{
				Name:        "plan_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Flavors").Transform(extractPlanName),
				Description: "Name of the savings plan.",
			},
			// Raw data for advanced queries
			{
				Name:        "flavors",
				Type:        proto.ColumnType_JSON,
				Description: "Complete flavors data (JSON array).",
			},
		},
	}
}

type CloudProjectUsagePlan struct {
	ProjectID    string        `json:"projectId"`
	Period       Period        `json:"period"`
	TotalSavings PlanPrice     `json:"totalSavings"`
	Flavors      []FlavorUsage `json:"flavors"`
}

type Period struct {
	From *time.Time `json:"from"`
	To   *time.Time `json:"to"`
}

type PlanPrice struct {
	CurrencyCode  string  `json:"currencyCode"`
	PriceInUcents int64   `json:"priceInUcents"`
	Text          string  `json:"text"`
	Value         float64 `json:"value"`
}

type FlavorUsage struct {
	Flavor        string             `json:"flavor"`
	Fees          FlavorFees         `json:"fees"`
	Periods       []UsagePeriod      `json:"periods"`
	Subscriptions []PlanSubscription `json:"subscriptions"`
}

type FlavorFees struct {
	FlatFee     FlatFee   `json:"flatFee"`
	OverQuota   OverQuota `json:"overQuota"`
	SavedAmount PlanPrice `json:"savedAmount"`
	TotalPrice  PlanPrice `json:"totalPrice"`
}

type FlatFee struct {
	Details    []FlatFeeDetail `json:"details"`
	TotalPrice PlanPrice       `json:"totalPrice"`
}

type FlatFeeDetail struct {
	ID         string    `json:"id"`
	Period     Period    `json:"period"`
	PlanName   string    `json:"planName"`
	Size       int       `json:"size"`
	TotalPrice PlanPrice `json:"totalPrice"`
	UnitPrice  PlanPrice `json:"unitPrice"`
}

type OverQuota struct {
	IDs        []string  `json:"ids"`
	Quantity   int       `json:"quantity"`
	TotalPrice PlanPrice `json:"totalPrice"`
	UnitPrice  PlanPrice `json:"unitPrice"`
}

type UsagePeriod struct {
	PlansIDs        []string   `json:"plansIds"`
	Begin           *time.Time `json:"begin"`
	ConsumptionSize int        `json:"consumptionSize"`
	Coverage        string     `json:"coverage"`
	CumulPlanSize   int        `json:"cumulPlanSize"`
	End             *time.Time `json:"end"`
	ResourceIDs     []string   `json:"resourceIds"`
	Utilization     string     `json:"utilization"`
}

type PlanSubscription struct {
	Begin *time.Time `json:"begin"`
	End   *time.Time `json:"end"`
	ID    string     `json:"id"`
	Size  int        `json:"size"`
}

func listCloudProjectUsagePlans(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	projectID := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCloudProjectUsagePlans.connect", "connection_error", err)
		return nil, err
	}

	endpoint := fmt.Sprintf("/cloud/project/%s/usage/plans", projectID)
	var response CloudProjectUsagePlan

	if err := client.Get(endpoint, &response); err != nil {
		plugin.Logger(ctx).Error("listCloudProjectUsagePlans", "api_error", err)
		return nil, err
	}

	// Stream the response as a single item
	d.StreamListItem(ctx, response)

	return nil, nil
}

// Helper functions to extract data from the first flavor in the flavors array

func extractFirstFlavorName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 {
		return nil, nil
	}

	return flavors[0].Flavor, nil
}

func extractFlatFeeTotalPrice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 {
		return nil, nil
	}

	return flavors[0].Fees.FlatFee.TotalPrice.Value, nil
}

func extractFlatFeeCurrency(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 {
		return nil, nil
	}

	return flavors[0].Fees.FlatFee.TotalPrice.CurrencyCode, nil
}

func extractOverQuotaQuantity(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 {
		return nil, nil
	}

	return flavors[0].Fees.OverQuota.Quantity, nil
}

func extractOverQuotaUnitPrice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 {
		return nil, nil
	}

	return flavors[0].Fees.OverQuota.UnitPrice.Value, nil
}

func extractFlavorTotalPrice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 {
		return nil, nil
	}

	return flavors[0].Fees.TotalPrice.Value, nil
}

func extractFlavorSavedAmount(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 {
		return nil, nil
	}

	return flavors[0].Fees.SavedAmount.Value, nil
}

func extractUsageCoverage(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 || len(flavors[0].Periods) == 0 {
		return nil, nil
	}

	return flavors[0].Periods[0].Coverage, nil
}

func extractUsageUtilization(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 || len(flavors[0].Periods) == 0 {
		return nil, nil
	}

	return flavors[0].Periods[0].Utilization, nil
}

func extractConsumptionSize(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 || len(flavors[0].Periods) == 0 {
		return nil, nil
	}

	return flavors[0].Periods[0].ConsumptionSize, nil
}

func extractCumulPlanSize(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 || len(flavors[0].Periods) == 0 {
		return nil, nil
	}

	return flavors[0].Periods[0].CumulPlanSize, nil
}

func extractSubscriptionID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 || len(flavors[0].Subscriptions) == 0 {
		return nil, nil
	}

	return flavors[0].Subscriptions[0].ID, nil
}

func extractSubscriptionSize(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 || len(flavors[0].Subscriptions) == 0 {
		return nil, nil
	}

	return flavors[0].Subscriptions[0].Size, nil
}

func extractSubscriptionBegin(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 || len(flavors[0].Subscriptions) == 0 {
		return nil, nil
	}

	return flavors[0].Subscriptions[0].Begin, nil
}

func extractSubscriptionEnd(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 || len(flavors[0].Subscriptions) == 0 {
		return nil, nil
	}

	return flavors[0].Subscriptions[0].End, nil
}

func extractPlanName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	flavors, ok := d.Value.([]FlavorUsage)
	if !ok || len(flavors) == 0 || len(flavors[0].Fees.FlatFee.Details) == 0 {
		return nil, nil
	}

	return flavors[0].Fees.FlatFee.Details[0].PlanName, nil
}
