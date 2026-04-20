# Table: ovh_deposit_paid_bill

Paid bills associated with deposits in your OVH account.

The `ovh_deposit_paid_bill` table can be used to query information about bills that were paid using account deposits. This table supports efficient querying by deposit ID for fast queries, or by deposit date for range-based queries.

## Examples

### List all paid bills for a specific deposit (fast path)

```sql
select
  deposit_id,
  bill_id,
  price_with_tax_value,
  price_with_tax_currency,
  payment_date
from
  ovh_deposit_paid_bill
where
  deposit_id = 'PA_FRxxxxxxxx'
order by
  payment_date desc;
```

### Get paid bills from a specific month (smart enumeration)

```sql
select
  deposit_id,
  deposit_date,
  bill_id,
  price_with_tax_value,
  category,
  payment_date
from
  ovh_deposit_paid_bill
where
  deposit_date >= '2024-01-01'::timestamp
  and deposit_date <= '2024-01-31'::timestamp
order by
  deposit_date desc;
```

### Get a specific bill for a specific deposit

```sql
select
  deposit_id,
  bill_id,
  date,
  price_with_tax_value,
  price_without_tax_value,
  tax_value,
  category,
  pdf_url,
  payment_type,
  payment_date
from
  ovh_deposit_paid_bill
where
  deposit_id = 'PA_FRxxxxxxxx'
  and bill_id = 'FRxxxxxxxx';
```

### Get all paid bills with deposit and bill dates

```sql
select
  deposit_id,
  deposit_date,
  bill_id,
  date as bill_date,
  price_with_tax_value,
  payment_identifier,
  payment_date
from
  ovh_deposit_paid_bill
where
  deposit_date >= '2024-01-01'::timestamp
  and deposit_date <= '2024-12-31'::timestamp
order by
  deposit_date desc,
  bill_id;
```

### Join deposits with their paid bills to build a complete lineage

```sql
select
  d.deposit_id,
  d.date as deposit_date,
  d.amount_value as deposit_amount,
  b.bill_id,
  b.date as bill_date,
  b.price_with_tax_value as bill_amount,
  b.category,
  b.payment_type,
  b.payment_date
from
  ovh_deposit d
left join ovh_deposit_paid_bill b on d.deposit_id = b.deposit_id
where
  d.date >= '2024-01-01'::timestamp
  and d.date <= '2024-12-31'::timestamp
order by
  d.date desc,
  b.bill_id;
```

### Get invoicing status for all paid bills in a date range

```sql
select
  deposit_id,
  deposit_date,
  bill_id,
  date,
  category,
  e_invoicing_id,
  e_invoicing_status,
  pdf_url
from
  ovh_deposit_paid_bill
where
  deposit_date >= '2024-06-01'::timestamp
  and deposit_date <= '2024-06-30'::timestamp
order by
  date desc;
```

## Performance Notes

The table uses intelligent query optimization:

- **Fast Path** (with `deposit_id`): Returns results in milliseconds (~2-3 API calls)
- **Smart Path** (with `deposit_date` range): Returns results in seconds (~12-50 API calls for 1 month)
- **Full Enumeration** (no filters): Slower performance (~547+ API calls for all deposits)

For best performance, always include a `deposit_id` qualifier when possible, or use `deposit_date` range filters to limit the scope of enumeration.
