# Table: ovh_cloud_data_job

A data job is processed by OVH by Apache Spark.

The `ovh_cloud_data_job` table can be used to query information about your jobs and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List data job of a cloud project

```sql
select
  id,
  name,
  status,
from
  ovh_cloud_data_job
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```

### List completed data job

```sql
select
  id,
  name,
  status,
from
  ovh_cloud_data_job
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and status = 'COMPLETED'
```
