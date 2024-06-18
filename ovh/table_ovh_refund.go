package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type Refund struct {
	ID              string    `json:"refundId"`
	Date            time.Time `json:"date"`
	Url             string    `json:"url"`
	PdfUrl          string    `json:"pdfUrl"`
	OrderId         int       `json:"orderId"`
	OriginalBillId  string    `json:"originalBillId"`
	Password        string    `json:"password"`
	PriceWithTax    Price     `json:"priceWithTax"`
	PriceWithoutTax Price     `json:"priceWithoutTax"`
	Tax             Price     `json:"tax"`
}

func tableOvhRefund() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_refund",
		Description: "Refunds of your account.",
		List: &plugin.ListConfig{
			Hydrate: listRefund,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id"}),
			Hydrate:    getRefund,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{Func: getRefundInfo},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the refund.",
			},
			{
				Name:        "date",
				Hydrate:     getRefundInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date of the refund.",
			},
			{
				Name:        "url",
				Hydrate:     getRefundInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Url"),
				Description: "URL to download the refund document.",
			},
			{
				Name:        "pdf_url",
				Hydrate:     getRefundInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PdfUrl"),
				Description: "URL to download the refund document in PDF format (maybe same as url field).",
			},
			{
				Name:        "order_id",
				Hydrate:     getRefundInfo,
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("OrderId"),
				Description: "Order id.",
			},
			{
				Name:        "original_bill_id",
				Hydrate:     getRefundInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OriginalBillId"),
				Description: "Original Bill id.",
			},
			{
				Name:        "password",
				Hydrate:     getRefundInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Password to download the refund document.",
			},
			{
				Name:        "price_with_tax",
				Hydrate:     getRefundInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("PriceWithTax.Value"),
				Description: "Price with tax.",
			},
			{
				Name:        "price_without_tax",
				Hydrate:     getRefundInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("PriceWithoutTax.Value"),
				Description: "Price without tax.",
			},
			{
				Name:        "tax",
				Hydrate:     getRefundInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Tax.Value"),
				Description: "Amount of the tax.",
			},
		},
	}
}

func getRefundInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	refund := h.Item.(Refund)

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_refund.getRefundInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/me/refund/%s", refund.ID), &refund)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill.getRefundInfo", err)
		return nil, err
	}

	return refund, nil
}

func listRefund(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_refund.listRefund", "connection_error", err)
		return nil, err
	}

	var refundsId []string
	err = client.Get("/me/refund", &refundsId)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_refund.listRefund", err)
		return nil, err
	}

	for _, refundId := range refundsId {
		var refund Refund
		refund.ID = refundId
		d.StreamListItem(ctx, refund)
	}

	return nil, nil
}

func getRefund(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.Quals.ToEqualsQualValueMap()["id"].GetStringValue()
	var refund Refund
	refund.ID = id
	return refund, nil
}
