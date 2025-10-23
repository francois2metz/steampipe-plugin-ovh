# Table: ovh_cloud_project

A cloud project is a way to regroup instance, storage, database, ... under a name.

The `ovh_cloud_project` table can be used to query information about cloud projects.

## Examples

### List projects

```sql
select
  id,
  description
from
  ovh_cloud_project
```

### List projects with IAM metadata

```sql
select
  id,
  name,
  description,
  iam ->> 'id' as iam_id,
  iam ->> 'urn' as iam_urn,
  iam ->> 'displayName' as iam_display_name,
  iam -> 'tags' as iam_tags
from
  ovh_cloud_project
```

### Query IAM URN for a specific project

```sql
select
  id,
  name,
  iam ->> 'urn' as iam_urn
from
  ovh_cloud_project
where
  id = 'f5346fda646e4364a788199e8240e720'
```

### List projects with specific IAM tags

```sql
select
  id,
  name,
  iam ->> 'displayName' as iam_display_name,
  iam -> 'tags' as iam_tags
from
  ovh_cloud_project
where
  iam -> 'tags' is not null
```
