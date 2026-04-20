# Table: ovh_deposit

Deposits of your OVH account.

The `ovh_deposit` table can be used to query information about your account deposits. Deposits are prepayments made to your OVH account that can be used to pay bills.

## Examples

### List all deposits

```sql
select
  deposit_id,
  date,
  amount_value,
  amount_currency,
  order_id
from
  ovh_deposit;
```

### Get deposits from a specific month

```sql
select
  deposit_id,
  date,
  amount_value,
  amount_currency
from
  ovh_deposit
where
  date >= '2024-01-01'::timestamp
  and date <= '2024-01-31'::timestamp
order by
  date desc;
```

### Get a deposit by ID

```sql
select
  deposit_id,
  date,
  amount_value,
  url,
  pdf_url
from
  ovh_deposit
where
  deposit_id = 'PA_FRxxxxxxxx';
```

### Get deposits by order ID

```sql
select
  deposit_id,
  date,
  amount_value,
  order_id
from
  ovh_deposit
where
  order_id = 12345678;
```

### Get deposits with their associated paid bills

```sql
select
  d.deposit_id,
  d.date as deposit_date,
  d.amount_value as deposit_amount,
  b.bill_id,
  b.price_with_tax_value as bill_amount,
  b.payment_date
from
  ovh_deposit d
left join ovh_deposit_paid_bill b on d.deposit_id = b.deposit_id
where
  d.date >= '2024-01-01'::timestamp
  and d.date <= '2024-01-31'::timestamp
order by
  d.date desc;
```
