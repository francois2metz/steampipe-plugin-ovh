# Table: ovh_cloud_ai_notebook

OVHcloud AI Notebook gives a quick and simple start launching your Jupyter or VS Code notebooks in the cloud. See the [official guide](https://www.ovhcloud.com/en/public-cloud/ai-notebooks/).

The `ovh_cloud_ai_notebook` table can be used to query information about your AI notebooks and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List AI notebooks of a cloud project

```sql
select
  id,
  name,
  state
from
  ovh_cloud_ai_notebook
where
  project_id = '27c5a6d3dfez87893jfd88fdsfmvnqb8';
```

### List stopped AI notebooks

```sql
select
  id,
  name,
  state
from
  ovh_cloud_ai_notebook
where
  project_id = '27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and state = 'STOPPED';
```

### List AI notebooks using a specific framework

```sql
select
  id,
  name,
  framework,
  version,
  state
from
  ovh_cloud_ai_notebook
where
  project_id = '27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and framework = 'conda';
```
