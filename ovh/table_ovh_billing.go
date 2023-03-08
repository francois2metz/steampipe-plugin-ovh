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

type PriceApi struct {
	Value        float64 `json:"value"`
	CurrencyCode string  `json:"currencyCode"`
}

type BillingApi struct {
	ID              string    `json:"billId"`
	Date            time.Time `json:"date"`
	Url             string    `json:"url"`
	PdfUrl          string    `json:"pdfUrl"`
	OrderId         int       `json:"orderId"`
	Category        string    `json:"category"`
	Password        string    `json:"password"`
	PriceWithTax    PriceApi  `json:"priceWithTax"`
	PriceWithoutTax PriceApi  `json:"priceWithoutTax"`
	Tax             PriceApi  `json:"tax"`
}

func tableOvhBilling() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_billing",
		Description: "Billing of you account.",
		List: &plugin.ListConfig{
			Hydrate: listBilling,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id"}),
			Hydrate:    getBilling,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "ID of billing."},
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Description: "Date of billing."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Url"), Description: "URL to download billing."},
			{Name: "pdf_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("PdfUrl"), Description: "URL to download billing in PDF format (maybe same as url field)."},
			{Name: "order_id", Type: proto.ColumnType_INT, Transform: transform.FromField("OrderId"), Description: "Order id."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Category of billing (autorenew, earlyrenewal...)."},
			{Name: "password", Type: proto.ColumnType_STRING, Description: "Password to download billing."},
			{Name: "price_without_tax", Type: proto.ColumnType_DOUBLE, Transform: transform.FromField("PriceWithoutTax.Value"), Description: "Password to download billing."},
			{Name: "tax", Type: proto.ColumnType_DOUBLE, Transform: transform.FromField("Tax.Value"), Description: "Password to download billing."},
		},
	}
}

func getOneBill(ctx context.Context, client *ovh.Client, billId string) (BillingApi, error) {
	logger := plugin.Logger(ctx)

	if logger.IsDebug() {
		logger.Debug("ovh_billing.getOneBill", fmt.Sprintf("Get bill (id: %s)", billId))
	}

	var billApi BillingApi

	err := client.Get(fmt.Sprintf("/me/bill/%s", billId), &billApi)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_billing.getOneBill", err)
		return BillingApi{}, err
	}

	if logger.IsDebug() {
		logger.Debug("ovh_billing.getOneBill", fmt.Sprintf("Bill api %v+", billApi))
	}

	return billApi, nil
}

func listBilling(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_billing.listBilling", "connection_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Debug("ovh_billing.listBilling", "Get list of billing")

	// First, we get IDs of billing
	var billingsId []string
	err = client.Get("/me/bill", &billingsId)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_billing.listBilling", err)
		return nil, err
	}

	for _, id := range billingsId {
		bill, errBill := getOneBill(ctx, client, id)

		if errBill != nil {
			return nil, err
		}

		d.StreamListItem(ctx, bill)
	}

	return nil, nil
}

func getBilling(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_billing.getBilling", "connection_error", err)
		return nil, err
	}

	id := d.Quals.ToEqualsQualValueMap()["id"].GetStringValue()

	billing, err := getOneBill(ctx, client, id)

	return billing, err
}
