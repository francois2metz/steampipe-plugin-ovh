# Table: ovh_cloud_ai_job

OVHcloud AI Training lets you train your AI, machine learning and deep learning models efficiently and easily, and optimise your GPU usage. See the [official guide](https://www.ovhcloud.com/en/public-cloud/ai-training/).

The `ovh_cloud_ai_job` table can be used to query information about your AI jobs and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List AI jobs of a cloud project

```sql
select
  id,
  name,
  state
from
  ovh_cloud_ai_job
where
  project_id = '27c5a6d3dfez87893jfd88fdsfmvnqb8';
```

### List completed AI jobs

```sql
select
  id,
  name,
  state
from
  ovh_cloud_ai_job
where
  project_id = '27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and state = 'DONE';
```

### List AI jobs using a specific image

```sql
select
  id,
  name,
  image,
  state
from
  ovh_cloud_ai_job
where
  project_id = '27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and image LIKE 'pytorch%';
```
