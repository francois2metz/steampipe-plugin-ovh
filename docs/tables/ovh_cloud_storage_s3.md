# Table: ovh_cloud_storage_s3

An S3 storage is an S3 object storage.

The `ovh_cloud_storage_s3` table can be used to query information about storage containers. You must specify which cloud project in the where clause (`where project_id=xxxx`). The region can optionally be specified in the where clause to filter results, or you can configure which regions to query using the `regions` parameter in your connection configuration.

**Note:** By default, this table will query all available regions. To improve performance, you can:
- Specify a region in your query's where clause: `where region='GRA'`
- Configure the `regions` parameter in your connection to limit which regions are queried globally

## Examples

### List S3 storage containers of a cloud project in all configured regions

```sql
select
  name,
  region,
  owner_id,
  objects_count,
  objects_size
from
  ovh_cloud_storage_s3
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```

### List S3 storage containers in a specific region

```sql
select
  name,
  owner_id,
  objects_count,
  objects_size
from
  ovh_cloud_storage_s3
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and region='GRA'
```

### List specific storage container

```sql
select
  name,
  owner_id,
  objects_count,
  objects_size
from
  ovh_cloud_storage_s3
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and region='GRA'
  and name='databucket'
```

### List S3 storage containers across multiple regions

```sql
select
  name,
  region,
  owner_id,
  objects_count,
  objects_size
from
  ovh_cloud_storage_s3
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and region in ('GRA', 'SBG', 'BHS')
```
