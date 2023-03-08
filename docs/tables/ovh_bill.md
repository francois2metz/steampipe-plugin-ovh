# Table: ovh_bill

List of all the bills of your account.

The `ovh_bill` table can be used to query information about your billing information.

## Examples

### List bills

```sql
select
  id,
  date,
  price_with_tax
from
  ovh_bill;
```

### Get a bill

```sql
select
  id,
  date,
  price_with_tax
from
  ovh_bill
where
  id = 'FRxxxxxxxx';
```
