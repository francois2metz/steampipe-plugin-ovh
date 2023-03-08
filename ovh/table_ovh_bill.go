package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type Price struct {
	Value        float64 `json:"value"`
	CurrencyCode string  `json:"currencyCode"`
}

type Bill struct {
	ID              string    `json:"billId"`
	Date            time.Time `json:"date"`
	Url             string    `json:"url"`
	PdfUrl          string    `json:"pdfUrl"`
	OrderId         int       `json:"orderId"`
	Category        string    `json:"category"`
	Password        string    `json:"password"`
	PriceWithTax    Price     `json:"priceWithTax"`
	PriceWithoutTax Price     `json:"priceWithoutTax"`
	Tax             Price     `json:"tax"`
}

func tableOvhBill() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_bill",
		Description: "Bills of your account.",
		List: &plugin.ListConfig{
			Hydrate: listBill,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id"}),
			Hydrate:    getBill,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{Func: getBillInfo},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the bill.",
			},
			{
				Name:        "date",
				Hydrate:     getBillInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date of the bill.",
			},
			{
				Name:        "url",
				Hydrate:     getBillInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Url"),
				Description: "URL to download the bill.",
			},
			{
				Name:        "pdf_url",
				Hydrate:     getBillInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PdfUrl"),
				Description: "URL to download the bill in PDF format (maybe same as url field).",
			},
			{
				Name:        "order_id",
				Hydrate:     getBillInfo,
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("OrderId"),
				Description: "Order id.",
			},
			{
				Name:        "category",
				Hydrate:     getBillInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Category of the bill (autorenew, earlyrenewal...).",
			},
			{
				Name:        "password",
				Hydrate:     getBillInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Password to download the bill.",
			},
			{
				Name:        "price_with_tax",
				Hydrate:     getBillInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("PriceWithTax.Value"),
				Description: "Price with tax.",
			},
			{
				Name:        "price_without_tax",
				Hydrate:     getBillInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("PriceWithoutTax.Value"),
				Description: "Price without tax.",
			},
			{
				Name:        "tax",
				Hydrate:     getBillInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Tax.Value"),
				Description: "Amount of the tax.",
			},
		},
	}
}

func getBillInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bill := h.Item.(Bill)

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill.getBillInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/me/bill/%s", bill.ID), &bill)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill.getBillInfo", err)
		return nil, err
	}

	return bill, nil
}

func listBill(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill.listBill", "connection_error", err)
		return nil, err
	}

	var billsId []string
	err = client.Get("/me/bill", &billsId)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill.listBill", err)
		return nil, err
	}

	for _, billId := range billsId {
		var bill Bill
		bill.ID = billId
		d.StreamListItem(ctx, bill)
	}

	return nil, nil
}

func getBill(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.Quals.ToEqualsQualValueMap()["id"].GetStringValue()
	var bill Bill
	bill.ID = id
	return bill, nil
}
