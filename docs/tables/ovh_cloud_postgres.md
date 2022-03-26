# Table: ovh_cloud_postgres

An hosted postgres database.

The `ovh_cloud_postgres` table can be used to query information about postgres and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List postgres instances of a cloud project

```sql
select
  id,
  plan,
  status,
  description
from
  ovh_cloud_postgres
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```
