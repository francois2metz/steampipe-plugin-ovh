# Table: ovh_iam_resource

IAM resources represent all resources in your OVH account that can be managed through Identity and Access Management (IAM) policies. This includes dedicated servers, public cloud projects, IP addresses, vRacks, and more.

The `ovh_iam_resource` table can be used to query information about all IAM resources in your account.

## Examples

### List all IAM resources

```sql
select
  id,
  name,
  display_name,
  type,
  owner
from
  ovh_iam_resource
order by
  owner,
  type
```

### List resources by type

```sql
select
  id,
  name,
  display_name,
  type
from
  ovh_iam_resource
where
  type = 'publicCloudProject'
```

### Count resources by type

```sql
select
  type,
  count(*) as resource_count
from
  ovh_iam_resource
group by
  type
order by
  resource_count desc
```

### List all dedicated servers

```sql
select
  id,
  name,
  display_name,
  owner
from
  ovh_iam_resource
where
  type = 'dedicatedServer'
order by
  name
```

### List all public cloud projects

```sql
select
  id,
  name,
  display_name,
  owner
from
  ovh_iam_resource
where
  type = 'publicCloudProject'
```

### List IP addresses with their routing information

```sql
select
  name as ip_address,
  display_name,
  tags ->> 'ovh:routedTo' as routed_to,
  tags ->> 'ovh:type' as ip_type,
  tags ->> 'ovh:campus' as campus,
  tags ->> 'ovh:isAdditionalIp' as is_additional_ip
from
  ovh_iam_resource
where
  type = 'ip'
  and tags is not null
order by
  ip_address
```

### List failover IPs

```sql
select
  name as ip_address,
  display_name,
  tags ->> 'ovh:routedTo' as routed_to,
  tags ->> 'ovh:campus' as campus
from
  ovh_iam_resource
where
  type = 'ip'
  and tags ->> 'ovh:type' = 'failover'
```

### List cloud IPs routed to a specific project

```sql
select
  name as ip_address,
  tags ->> 'ovh:routedTo' as project_id,
  tags ->> 'ovh:version' as ip_version
from
  ovh_iam_resource
where
  type = 'ip'
  and tags ->> 'ovh:type' = 'cloud'
  and tags ->> 'ovh:routedTo' = 'your-project-id'
```

### List resources owned by a specific account

```sql
select
  type,
  count(*) as count
from
  ovh_iam_resource
where
  owner = 'your-account-ovh'
group by
  type
order by
  count desc
```

### List vRacks

```sql
select
  id,
  name,
  display_name,
  owner
from
  ovh_iam_resource
where
  type = 'vrack'
```

### List resources with tags

```sql
select
  id,
  name,
  type,
  display_name,
  jsonb_pretty(tags) as tags
from
  ovh_iam_resource
where
  tags is not null
  and jsonb_typeof(tags) = 'object'
limit 10
```

### Find resources by campus location

```sql
select
  name,
  type,
  display_name,
  tags ->> 'ovh:campus' as campus
from
  ovh_iam_resource
where
  tags ->> 'ovh:campus' = 'GRA'
```

### List IPv6 addresses

```sql
select
  name as ip_address,
  display_name,
  tags ->> 'ovh:routedTo' as routed_to,
  tags ->> 'ovh:campus' as campus
from
  ovh_iam_resource
where
  type = 'ip'
  and tags ->> 'ovh:version' = '6'
order by
  ip_address
```

### List additional IPs

```sql
select
  name as ip_address,
  display_name,
  tags ->> 'ovh:routedTo' as routed_to,
  tags ->> 'ovh:type' as ip_type
from
  ovh_iam_resource
where
  type = 'ip'
  and tags ->> 'ovh:isAdditionalIp' = 'true'
