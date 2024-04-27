# Table: ovh_cloud_s3_storage

An S3 storage is an S3 object storage.

The `ovh_cloud_s3_storage` table can be used to query information about storage containers and **you must specify which cloud project AND region** in the where clause (`where project_id=xxxx and region=xxxx`).

## Examples

### List S3 storage containers of a cloud project

```sql
select
  name,
  owner_id,
  objects_count,
  objects_size
from
  ovh_cloud_s3_storage
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and region='GRA'
```

## List specific storage container

```sql
select
  name,
  owner_id,
  objects_count,
  objects_size
from
  ovh_cloud_s3_storage
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and region='GRA'
  and name='databucket'
```
