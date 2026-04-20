package ovh

import (
	"context"
	"fmt"
	"net/url"
	"time"

	ovhapi "github.com/ovh/go-ovh/ovh"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type DepositPaidBill struct {
	DepositID         string     `json:"deposit_id"`
	DepositDate       time.Time  `json:"depositDate"`
	BillID            string     `json:"billId"`
	OrderID           int        `json:"orderId"`
	PriceWithTax      Price      `json:"priceWithTax"`
	PriceWithoutTax   Price      `json:"priceWithoutTax"`
	Tax               Price      `json:"tax"`
	Date              time.Time  `json:"date"`
	Category          string     `json:"category"`
	Url               string     `json:"url"`
	PdfUrl            string     `json:"pdfUrl"`
	Password          string     `json:"password"`
	EInvoicingID      *string    `json:"eInvoicingId"`
	EInvoicingStatus  *string    `json:"eInvoicingStatus"`
	PaymentType       *string    `json:"paymentType"`
	PaymentIdentifier *string    `json:"paymentIdentifier"`
	PaymentDate       *time.Time `json:"paymentDate"`
}

func tableOvhDepositPaidBill() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_deposit_paid_bill",
		Description: "Paid bills associated with deposits. Supports querying by deposit_id or deposit_date (or both).",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "deposit_id", Operators: []string{"="}, Require: plugin.Optional},
				{Name: "deposit_date", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
			},
			Hydrate: listDepositPaidBills,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"deposit_id", "bill_id"}),
			Hydrate:    getDepositPaidBill,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{Func: getDepositPaidBillInfo},
		},
		Columns: []*plugin.Column{
			{
				Name:        "deposit_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DepositID"),
				Description: "ID of the deposit.",
			},
			{
				Name:        "deposit_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DepositDate"),
				Description: "Date of the deposit (for easy filtering).",
			},
			{
				Name:        "bill_id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the paid bill.",
			},
			{
				Name:        "order_id",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("OrderID"),
				Description: "Order ID associated with the bill.",
			},
			{
				Name:        "price_with_tax_value",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("PriceWithTax.Value"),
				Description: "Price with tax.",
			},
			{
				Name:        "price_with_tax_currency",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PriceWithTax.CurrencyCode"),
				Description: "Currency code for price with tax.",
			},
			{
				Name:        "price_without_tax_value",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("PriceWithoutTax.Value"),
				Description: "Price without tax.",
			},
			{
				Name:        "price_without_tax_currency",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PriceWithoutTax.CurrencyCode"),
				Description: "Currency code for price without tax.",
			},
			{
				Name:        "tax_value",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Tax.Value"),
				Description: "Tax amount.",
			},
			{
				Name:        "tax_currency",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Tax.CurrencyCode"),
				Description: "Currency code for tax.",
			},
			{
				Name:        "date",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date of the bill.",
			},
			{
				Name:        "category",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Category of the bill.",
			},
			{
				Name:        "url",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Url"),
				Description: "URL to access the bill document.",
			},
			{
				Name:        "pdf_url",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PdfUrl"),
				Description: "URL to download the bill document in PDF format.",
			},
			{
				Name:        "password",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Password to access the bill document.",
			},
			{
				Name:        "e_invoicing_id",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EInvoicingID"),
				Description: "E-invoicing identifier.",
			},
			{
				Name:        "e_invoicing_status",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EInvoicingStatus"),
				Description: "E-invoicing status.",
			},
			{
				Name:        "payment_type",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Type of payment.",
			},
			{
				Name:        "payment_identifier",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Payment identifier.",
			},
			{
				Name:        "payment_date",
				Hydrate:     getDepositPaidBillInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date of payment.",
			},
		},
	}
}

func getDepositPaidBillInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bill := h.Item.(DepositPaidBill)

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_deposit_paid_bill.getDepositPaidBillInfo", "connection_error", err)
		return nil, err
	}

	// Fetch bill details
	err = client.Get(fmt.Sprintf("/me/deposit/%s/paidBills/%s", bill.DepositID, bill.BillID), &bill)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_deposit_paid_bill.getDepositPaidBillInfo", "fetch_bill_error", err)
		return nil, err
	}

	// Fetch payment information
	var payment struct {
		PaymentType       *string    `json:"paymentType"`
		PaymentIdentifier *string    `json:"paymentIdentifier"`
		PaymentDate       *time.Time `json:"paymentDate"`
	}

	err = client.Get(fmt.Sprintf("/me/deposit/%s/paidBills/%s/payment", bill.DepositID, bill.BillID), &payment)
	if err != nil {
		// Log but don't fail - payment info might not exist
		plugin.Logger(ctx).Debug("ovh_deposit_paid_bill.getDepositPaidBillInfo", "fetch_payment_error", err)
	} else {
		bill.PaymentType = payment.PaymentType
		bill.PaymentIdentifier = payment.PaymentIdentifier
		bill.PaymentDate = payment.PaymentDate
	}

	return bill, nil
}

func listDepositPaidBills(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_deposit_paid_bill.listDepositPaidBills", "connection_error", err)
		return nil, err
	}

	depositID := d.EqualsQuals["deposit_id"].GetStringValue()

	// Path A: Fast path - if deposit_id is provided
	if depositID != "" {
		// Fetch deposit info to get the date
		var deposit Deposit
		err = client.Get(fmt.Sprintf("/me/deposit/%s", depositID), &deposit)
		if err != nil {
			plugin.Logger(ctx).Error("ovh_deposit_paid_bill.listDepositPaidBills", "fetch_deposit_error", err)
			return nil, err
		}

		var billIDs []string
		err = client.Get(fmt.Sprintf("/me/deposit/%s/paidBills", depositID), &billIDs)

		if err != nil {
			plugin.Logger(ctx).Error("ovh_deposit_paid_bill.listDepositPaidBills", "fetch_bills_error", err)
			return nil, err
		}

		for _, billID := range billIDs {
			var bill DepositPaidBill
			bill.DepositID = depositID
			bill.DepositDate = deposit.Date
			bill.BillID = billID
			d.StreamListItem(ctx, bill)
		}

		return nil, nil
	}

	// Path B: Smart path - enumerate all deposits (optionally filtered by date)
	deposits, err := enumerateDeposits(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_deposit_paid_bill.listDepositPaidBills", "enumerate_deposits_error", err)
		return nil, err
	}

	// For each deposit, fetch its paid bills
	for _, deposit := range deposits {
		var billIDs []string
		err = client.Get(fmt.Sprintf("/me/deposit/%s/paidBills", deposit.DepositID), &billIDs)

		if err != nil {
			plugin.Logger(ctx).Debug("ovh_deposit_paid_bill.listDepositPaidBills", "fetch_bills_error", fmt.Sprintf("deposit %s: %v", deposit.DepositID, err))
			continue // Skip this deposit and continue with next
		}

		// Stream each bill with its deposit date
		for _, billID := range billIDs {
			var bill DepositPaidBill
			bill.DepositID = deposit.DepositID
			bill.DepositDate = deposit.Date
			bill.BillID = billID
			d.StreamListItem(ctx, bill)
		}
	}

	return nil, nil
}

// enumerateDeposits fetches all deposits, optionally filtered by date range
func enumerateDeposits(ctx context.Context, d *plugin.QueryData, client *ovhapi.Client) ([]Deposit, error) {
	// Build query parameters for date filtering
	params := url.Values{}

	// Handle deposit_date filtering
	if d.Quals["deposit_date"] != nil {
		for _, qual := range d.Quals["deposit_date"].Quals {
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

	// Build the URL with query parameters
	depositURL := "/me/deposit"
	if len(params) > 0 {
		depositURL = depositURL + "?" + params.Encode()
	}

	var depositsID []string
	err := client.Get(depositURL, &depositsID)
	if err != nil {
		return nil, err
	}

	// Fetch full deposit details for each ID
	deposits := make([]Deposit, 0, len(depositsID))
	for _, depositID := range depositsID {
		var deposit Deposit
		deposit.DepositID = depositID

		// Fetch full deposit info
		err = client.Get(fmt.Sprintf("/me/deposit/%s", depositID), &deposit)
		if err != nil {
			plugin.Logger(ctx).Debug("ovh_deposit_paid_bill.enumerateDeposits", "fetch_deposit_error", fmt.Sprintf("deposit %s: %v", depositID, err))
			continue // Skip deposits that fail to fetch
		}

		deposits = append(deposits, deposit)
	}

	return deposits, nil
}

func getDepositPaidBill(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.Quals.ToEqualsQualValueMap()
	depositID := quals["deposit_id"].GetStringValue()
	billID := quals["bill_id"].GetStringValue()

	var bill DepositPaidBill
	bill.DepositID = depositID
	bill.BillID = billID
	return bill, nil
}
