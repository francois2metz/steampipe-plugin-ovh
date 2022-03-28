# Table: ovh_cloud_storage

A storage is an object storage similar to S3.

The `ovh_cloud_storage` table can be used to query information about storage containers and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List storage containers of a cloud project

```sql
select
  id,
  name
from
  ovh_cloud_storage
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```

## List empty storage containers

```sql
select
  id,
  name
from
  ovh_cloud_storage
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and stored_objects is null
```
