# Table: ovh_log_self

List of all the logs from your account.

The `ovh_log_self` table can be used to query information about your billing information.

## Examples

### List logs

```sql
select
  id,
  date,
  acount
from
  ovh_log_self;
```

### Get a log

```sql
select
  id,
  date,
  account
from
  ovh_log_self
where
  id = 'XXXXXX';
```
