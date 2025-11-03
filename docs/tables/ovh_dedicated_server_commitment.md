# Table: ovh_dedicated_server_commitment

Lists all dedicated servers and their commitment/renewal information. This table provides insights into server engagement periods, renewal settings, and service lifecycle information.

The `ovh_dedicated_server_commitment` table can be used to query information about dedicated server commitments and renewal configurations.

## Examples

### List all dedicated server commitments

```sql
select
  server_name,
  engaged_up_to,
  renew_mode,
  expiration
from
  ovh_dedicated_server_commitment;
```

### Get servers with automatic renewal enabled

```sql
select
  server_name,
  status,
  renew_mode,
  renew_period,
  engaged_up_to,
  expiration
from
  ovh_dedicated_server_commitment
where
  renew_mode = 'automatic';
```

### Find servers with engagement ending soon

```sql
select
  server_name,
  engaged_up_to,
  expiration,
  renew_mode
from
  ovh_dedicated_server_commitment
where
  engaged_up_to < now() + interval '30 days'
  and engaged_up_to is not null;
```

### Get detailed commitment information for a specific server

```sql
select
  server_name,
  service_id,
  status,
  renew_mode,
  renew_period,
  engaged_up_to,
  expiration,
  creation
from
  ovh_dedicated_server_commitment
where
  server_name = 'ns123456.ip-91-121-12.eu';
```

### List servers by renewal period

```sql
select
  renew_period,
  count(*) as server_count
from
  ovh_dedicated_server_commitment
group by
  renew_period
order by
  renew_period;
```

### Find servers without commitment

```sql
select
  server_name,
  status,
  renew_mode,
  expiration
from
  ovh_dedicated_server_commitment
where
  engaged_up_to is null;
```