package ovh

import (
	"context"
	"fmt"
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
	Start       string `json:"periodStart"`
	End         string `json:"periodEnd"`
	Quantity    string `json:"quantity"`
}

func tableOvhBillDetails() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_bill_details",
		Description: "Details' bill of you account.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"bill_id"}),
			Hydrate:    listBillingDetails,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"bill_id", "id"}),
			Hydrate:    getBillingDetail,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "ID of detail's bill."},
			{Name: "bill_id", Transform: transform.FromQual("bill_id"), Type: proto.ColumnType_STRING, Description: "ID of bill."},
			{Name: "description", Hydrate: getGetBillDetailInfo, Type: proto.ColumnType_STRING, Description: "Description of detail."},
			{Name: "domain", Hydrate: getGetBillDetailInfo, Type: proto.ColumnType_STRING, Description: "Domain."},
			{Name: "start", Transform: transform.From(convertBillDetailDate), Hydrate: getGetBillDetailInfo, Type: proto.ColumnType_TIMESTAMP, Description: "Start date of detail."},
			{Name: "end", Transform: transform.From(convertBillDetailDate), Hydrate: getGetBillDetailInfo, Type: proto.ColumnType_TIMESTAMP, Description: "End date of detail."},
			{Name: "quantity", Hydrate: getGetBillDetailInfo, Type: proto.ColumnType_STRING, Description: "Quantity of detail."},
		},
	}
}

func convertBillDetailDate(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	billDetail := d.HydrateItem.(BillDetail)

	var value string

	switch d.ColumnName {
	case "start":
		value = billDetail.Start
	case "end":
		value = billDetail.End
	}

	if len(value) == 0 {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02", value)

	return t, err
}

// This function populate data of bill detail
func getGetBillDetailInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	billDetail := h.Item.(BillDetail)

	client, err := connect(ctx, d)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill_details.getGetBillDetailInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/me/bill/%s/details/%s", billDetail.BillID, billDetail.ID), &billDetail)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_data_job.getDataJobInfo", err)
		return nil, err
	}

	return billDetail, nil
}

func listBillingDetails(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill_details.listBillingDetails", "connection_error", err)
		return nil, err
	}

	billId := d.EqualsQuals["bill_id"].GetStringValue()

	// First, we get IDs of billing
	var billDetailsId []string
	err = client.Get(fmt.Sprintf("/me/bill/%s/details", billId), &billDetailsId)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_bill_details.listBillingDetails", err)
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

	return getGetBillDetailInfo(ctx, d, h)
}
