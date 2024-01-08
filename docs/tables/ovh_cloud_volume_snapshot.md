# Table: ovh_cloud_volume_snapshot

A volume snapshot a copy of the state of a storage volume at a particular point in time.

The `ovh_cloud_volume_snapshot` table can be used to query information about volumes snapshots and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List volumes snapshots of a cloud project

```sql
select
  id,
  name
from
  ovh_cloud_volume_snapshot
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```

### List available volumes snapshots of a cloud project

```sql
select
  id,
  name
from
  ovh_cloud_volume_snapshot
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and status='available'
```

### Get one snapshot

```sql
select
  id,
  name
from
  ovh_cloud_volume_snapshot
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and id='d82f4336533f4d16ad92bf7b9f3d87e2'
```
