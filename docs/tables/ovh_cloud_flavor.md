# Table: ovh_cloud_flavor

A flavor is the instance model defining its characteristics in terms of resources.

The `ovh_cloud_flavor` table can be used to query information about flavors and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List flavor of a cloud project

```sql
select
  id,
  name,
  type
from
  ovh_cloud_flavor
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```
