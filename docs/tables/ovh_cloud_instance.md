# Table: ovh_cloud_instance

An instance is a virtual server in the OVH cloud.

The `ovh_cloud_instance` table can be used to query information about instances and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List instances of a cloud project

```sql
select
  id,
  name
from
  ovh_cloud_instance
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```

### List all instances of all cloud projects

```sql
select
  ci.id,
  ci.name
from
  ovh_cloud_instance ci
join
  ovh_cloud_project cp
on
  ci.project_id = cp.id
```
