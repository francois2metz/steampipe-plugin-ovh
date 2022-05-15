# Table: ovh_cloud_ssh_key

An ssh key allows you to connect to an instance.

The `ovh_cloud_ssh_key` table can be used to query information about ssh keys and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List ssh keys of a cloud project

```sql
select
  id,
  name,
  public_key
from
  ovh_cloud_ssh_key
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```
