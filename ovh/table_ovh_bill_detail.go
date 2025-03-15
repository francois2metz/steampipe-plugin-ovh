package ovh

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type BillDetail struct {
	ID          string `json:"id"`
	BillID      string `json:"bill_id"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	PeriodStart string `json:"periodStart"`
	PeriodEnd   string `json:"periodEnd"`
	Quantity    string `json:"quantity"`
	TotalPrice  Price  `json:"totalPrice"`
	UnitPrice   Price  `json:"unitPrice"`
}

func tableOvhBillDetails() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_bill_detail",
		Description: "Detail of a bill.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"bill_id"}),
			Hydrate:    listBillingDetails,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"bill_id", "id"}),
			Hydrate:    getBillingDetail,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of detail's bill.",
			},
			{
				Name:        "bill_id",
				Transform:   transform.FromQual("bill_id"),
				Type:        proto.ColumnType_STRING,
				Description: "ID of bill.",
			},
			{
				Name:        "description",
				Hydrate:     getBillDetailInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Description of detail.",
			},
			{
				Name:        "domain",
				Hydrate:     getBillDetailInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Domain.",
			},
			{
				Name:        "period_start",
				Transform:   transform.FromP(convertBillDetailDate, "PeriodStart"),
				Hydrate:     getBillDetailInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Period start of the product detail.",
			},
			{
				Name:        "period_end",
				Transform:   transform.FromP(convertBillDetailDate, "PeriodEnd"),
				Hydrate:     getBillDetailInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Period end of the product detail.",
			},
			{
				Name:        "quantity",
				Hydrate:     getBillDetailInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Quantity of detail.",
			},
			{
				Name:        "total_price",
				Hydrate:     getBillDetailInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("TotalPrice.Value"),
				Description: "Total price of this detail.",
			},
			{
				Name:        "unit_price",
				Hydrate:     getBillDetailInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("UnitPrice.Value"),
				Description: "Unit price of this detail.",
			},
		},
	}
}

func convertBillDetailDate(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	billDetail := d.HydrateItem.(BillDetail)
	columnName := d.Param.(string)
	billDetailReflect := reflect.ValueOf(billDetail)
	value := billDetailReflect.FieldByName(columnName).String()
	if len(value) == 0 {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", value)

	return t, err
}

// This function populate data of bill detail
func getBillDetailInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	billDetail := h.Item.(BillDetail)

	client, err := connect(ctx, d)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill_detail.getBillDetailInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/me/bill/%s/details/%s", billDetail.BillID, billDetail.ID), &billDetail)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill_detail.getBillDetailInfo", err)
		return nil, err
	}

	return billDetail, nil
}

func listBillingDetails(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill_detail.listBillingDetails", "connection_error", err)
		return nil, err
	}

	billId := d.EqualsQuals["bill_id"].GetStringValue()

	// First, we get IDs of billing
	var billDetailsId []string
	err = client.Get(fmt.Sprintf("/me/bill/%s/details", billId), &billDetailsId)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill_detail.listBillingDetails", err)
		return nil, err
	}

	for _, id := range billDetailsId {
		d.StreamListItem(ctx, BillDetail{
			ID:     id,
			BillID: billId,
		})
	}

	return nil, nil
}

func getBillingDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	billId := d.EqualsQuals["bill_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	h.Item = BillDetail{
		ID:     id,
		BillID: billId,
	}

	return getBillDetailInfo(ctx, d, h)
}
