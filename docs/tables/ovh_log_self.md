# Table: ovh_log_self

Get the logs from recent API calls made from your account.

The `ovh_log_self` table can be used to query information about your recent API calls.

## Examples

### List logs

```sql
select
  id,
  date,
  account
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
