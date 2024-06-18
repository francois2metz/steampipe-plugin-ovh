# Table: ovh_refund

List of all the refunds of your account.

The `ovh_refund` table can be used to query information about your refund information.

## Examples

### List refunds

```sql
select
  id,
  date,
  original_bill_id,
  price_with_tax
from
  ovh_refund;
```

### Get a bill

```sql
select
  id,
  date,
  original_bill_id,
  price_with_tax
from
  ovh_refund
where
  id = 'AFRxxxxxxx';
```
