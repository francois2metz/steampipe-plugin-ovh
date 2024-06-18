package ovh

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type RefundPrice struct {
	Value        float64 `json:"value"`
	CurrencyCode string  `json:"currencyCode"`
}
type RefundDetail struct {
	ID          string      `json:"refundDetailId"`
	RefundID    string      `json:"refundId"`
	Description string      `json:"description"`
	Domain      string      `json:"domain"`
	Quantity    string      `json:"quantity"`
	TotalPrice  RefundPrice `json:"totalPrice"`
	UnitPrice   RefundPrice `json:"unitPrice"`
}

func tableOvhRefundDetails() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_refund_detail",
		Description: "Detail of a refund.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"refund_id"}),
			Hydrate:    listRefundDetails,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"refund_id", "id"}),
			Hydrate:    getRefundDetail,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of detail's bill.",
			},
			{
				Name:        "refund_id",
				Transform:   transform.FromQual("refund_id"),
				Type:        proto.ColumnType_STRING,
				Description: "ID of refund detail.",
			},
			{
				Name:        "description",
				Hydrate:     getGetRefundDetailInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Description of detail.",
			},
			{
				Name:        "domain",
				Hydrate:     getGetRefundDetailInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Domain.",
			},
			{
				Name:        "quantity",
				Hydrate:     getGetRefundDetailInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Quantity of detail.",
			},
			{
				Name:        "total_price",
				Hydrate:     getGetRefundDetailInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("TotalPrice.Value"),
				Description: "Total price of this detail.",
			},
			{
				Name:        "unit_price",
				Hydrate:     getGetRefundDetailInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("UnitPrice.Value"),
				Description: "Unit price of this detail.",
			},
		},
	}
}

// This function populate data of refund detail
func getGetRefundDetailInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	refundDetail := h.Item.(RefundDetail)

	client, err := connect(ctx, d)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_refund_detail.getGetBillDetailInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/me/refund/%s/details/%s", refundDetail.RefundID, refundDetail.ID), &refundDetail)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_refund_detail.getGetBillDetailInfo", err)
		return nil, err
	}

	return refundDetail, nil
}

func listRefundDetails(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_refund_detail.listRefundDetails", "connection_error", err)
		return nil, err
	}

	refundId := d.EqualsQuals["refund_id"].GetStringValue()

	// First, we get IDs of refund
	var refundDetailsId []string
	err = client.Get(fmt.Sprintf("/me/refund/%s/details", refundId), &refundDetailsId)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_refund_detail.listRefundDetails", err)
		return nil, err
	}

	for _, id := range refundDetailsId {
		d.StreamListItem(ctx, RefundDetail{
			ID:       id,
			RefundID: refundId,
		})
	}

	return nil, nil
}

func getRefundDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	refundId := d.EqualsQuals["refund_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	h.Item = RefundDetail{
		ID:       id,
		RefundID: refundId,
	}

	return getGetRefundDetailInfo(ctx, d, h)
}
