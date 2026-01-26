package ovh

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type Deposit struct {
	DepositID   string     `json:"depositId"`
	Amount      Price      `json:"amount"`
	Date        time.Time  `json:"date"`
	OrderID     *int       `json:"orderId"`
	Url         string     `json:"url"`
	PdfUrl      string     `json:"pdfUrl"`
	Password    string     `json:"password"`
	PaymentInfo *string    `json:"paymentInfo"`
}

func tableOvhDeposit() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_deposit",
		Description: "Deposits of your account.",
		List: &plugin.ListConfig{
			Hydrate: listDeposit,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "date", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "order_id", Operators: []string{"="}, Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"deposit_id"}),
			Hydrate:    getDeposit,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{Func: getDepositInfo},
		},
		Columns: []*plugin.Column{
			{
				Name:        "deposit_id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the deposit.",
			},
			{
				Name:        "amount_value",
				Hydrate:     getDepositInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Amount.Value"),
				Description: "Deposit amount.",
			},
			{
				Name:        "amount_currency",
				Hydrate:     getDepositInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Amount.CurrencyCode"),
				Description: "Currency code of the deposit amount.",
			},
			{
				Name:        "date",
				Hydrate:     getDepositInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date of the deposit.",
			},
			{
				Name:        "order_id",
				Hydrate:     getDepositInfo,
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("OrderID"),
				Description: "Order ID associated with the deposit.",
			},
			{
				Name:        "url",
				Hydrate:     getDepositInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Url"),
				Description: "URL to access the deposit document.",
			},
			{
				Name:        "pdf_url",
				Hydrate:     getDepositInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PdfUrl"),
				Description: "URL to download the deposit document in PDF format.",
			},
			{
				Name:        "password",
				Hydrate:     getDepositInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Password to access the deposit document.",
			},
			{
				Name:        "payment_info",
				Hydrate:     getDepositInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Payment information for the deposit.",
			},
		},
	}
}

func getDepositInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	deposit := h.Item.(Deposit)

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_deposit.getDepositInfo", "connection_error", err)
		return nil, err
	}

	// Create a new deposit struct to hold the full response
	var fullDeposit Deposit
	
	// Try the /me/deposit/{id}/details endpoint
	err = client.Get(fmt.Sprintf("/me/deposit/%s/details", deposit.DepositID), &fullDeposit)

	if err != nil {
		// If /details endpoint fails, try without /details suffix
		plugin.Logger(ctx).Debug("ovh_deposit.getDepositInfo", "details_endpoint_error", err)
		err = client.Get(fmt.Sprintf("/me/deposit/%s", deposit.DepositID), &fullDeposit)
		if err != nil {
			plugin.Logger(ctx).Error("ovh_deposit.getDepositInfo", "fetch_error", err)
			return deposit, nil
		}
	}

	return fullDeposit, nil
}

func listDeposit(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_deposit.listDeposit", "connection_error", err)
		return nil, err
	}

	// Build query parameters
	params := url.Values{}

	// Handle date filtering
	if d.Quals["date"] != nil {
		for _, qual := range d.Quals["date"].Quals {
			if qual.Value != nil {
				dateValue := qual.Value.GetTimestampValue().AsTime().Format(time.RFC3339)
				switch qual.Operator {
				case ">=":
					params.Add("date.from", dateValue)
				case ">":
					dateValue = qual.Value.GetTimestampValue().AsTime().Add(time.Second).Format(time.RFC3339)
					params.Add("date.from", dateValue)
				case "<=":
					params.Add("date.to", dateValue)
				case "<":
					dateValue = qual.Value.GetTimestampValue().AsTime().Add(-time.Second).Format(time.RFC3339)
					params.Add("date.to", dateValue)
				case "=":
					params.Add("date.from", dateValue)
					params.Add("date.to", dateValue)
				}
			}
		}
	}

	// Handle order ID filtering
	if d.Quals["order_id"] != nil {
		orderID := d.Quals["order_id"].Quals[0].Value.GetInt64Value()
		params.Add("orderId", fmt.Sprintf("%d", orderID))
	}

	// Build the URL with query parameters
	depositURL := "/me/deposit"
	if len(params) > 0 {
		depositURL = depositURL + "?" + params.Encode()
	}

	var depositsID []string
	err = client.Get(depositURL, &depositsID)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_deposit.listDeposit", "fetch_error", err)
		return nil, err
	}

	for _, depositID := range depositsID {
		var deposit Deposit
		deposit.DepositID = depositID
		d.StreamListItem(ctx, deposit)
	}

	return nil, nil
}

func getDeposit(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.Quals.ToEqualsQualValueMap()["deposit_id"].GetStringValue()
	var deposit Deposit
	deposit.DepositID = id
	return deposit, nil
}
