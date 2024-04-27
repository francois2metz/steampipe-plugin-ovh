# Table: ovh_cloud_region

Regions available for a cloud project.

The `ovh_cloud_region` table can be used to query information about regions and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List regions of a cloud project

```sql
select
  name,
  type,
  status
from
  ovh_cloud_region
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```

### List regions not UP

```sql
select
  name
from
  ovh_cloud_postgres
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and status!='UP'
```

### Get specific region

```sql
select
  name,
  type,
  status
from
  ovh_cloud_region
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and name='GRA'
```
