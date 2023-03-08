# Table: ovh_bill_details

Details of the bill of your account.

The `ovh_bill_details` table can be used to query information about your billing information.

## Examples

### List details of one bill

```sql
select
  *
from
  ovh_bill_details
where
  bill_id = 'FRxxxxxx';
```

### Get one detail of one bill

```sql
select
  *
from
  ovh_bill_details
where
  bill_id = 'FRxxxxxxxx'
  and id = 'FRxxxxxxxx';
```
