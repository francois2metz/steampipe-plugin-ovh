# Table: ovh_cloud_volume

A volume is an independent additional disk.

The `ovh_cloud_volume` table can be used to query information about volumes and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List volumes of a cloud project

```sql
select
  id,
  name
from
  ovh_cloud_volume
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```

### List unused volumes of a cloud project

```sql
select
  id,
  name
from
  ovh_cloud_volume
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and status != 'in-use'
```
