# Table: ovh_refund_detail

Details of the refund of your account.

The `ovh_refund_detail` table can be used to query information about your refund information.

## Examples

### List details of one refund

```sql
select
  *
from
  ovh_refund_detail
where
  refund_id = 'AFRxxxxxxx';
```

### Get one detail of one refund

```sql
select
  *
from
  ovh_refund_detail
where
  refund_id = 'AFRxxxxxxxx'
  and id = 'AFRxxxxxxxx';
```
