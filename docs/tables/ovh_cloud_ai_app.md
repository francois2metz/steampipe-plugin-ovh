# Table: ovh_cloud_ai_app

An AI app is run by OVHcloud AI Deploy.

The `ovh_cloud_ai_app` table can be used to query information about your AI apps and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List AI apps of a cloud project

```sql
select
  id,
  name,
  state
from
  ovh_cloud_ai_app
where
  project_id = '27c5a6d3dfez87893jfd88fdsfmvnqb8';
```

### List running AI apps

```sql
select
  id,
  name,
  state
from
  ovh_cloud_ai_app
where
  project_id = '27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and state = 'RUNNING';
```

### List AI apps using a specific image

```sql
`select
  id,
  name,
  image,
  state
from
  ovh_cloud_ai_app
where
  project_id = '27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and image LIKE 'pytorch%';`
```
