# Table: ovh_cloud_ai_app

OVHcloud AI Deploy lets you easily deploy machine learning models and applications to production, create your API access points effortlessly, and make effective predictions. See the [official guide](https://www.ovhcloud.com/en/public-cloud/ai-deploy/).

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
