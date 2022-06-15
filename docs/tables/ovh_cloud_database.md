# Table: ovh_cloud_database

An hosted database.

The `ovh_cloud_database` table can be used to query information about databases and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List database instances of a cloud project

```sql
select
  id,
  plan,
  status,
  description
from
  ovh_cloud_database
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```

### List not ready database instances of a cloud project

```sql
select
  id,
  plan,
  status,
  description
from
  ovh_cloud_database
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and status!='READY'
```

### List PostgreSQL database of a cloud project

```sql
select
  id,
  plan,
  status,
  description
from
  ovh_cloud_database
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and engine='postgresql'
```
