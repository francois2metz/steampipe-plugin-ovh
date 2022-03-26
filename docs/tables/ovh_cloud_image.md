# Table: ovh_cloud_image

An image is a pre-installed, ready-to-use operating system. 

The `ovh_cloud_image` table can be used to query information about images and **you must specify which cloud project** in the where or join clause (`where project_id=`, `join ovh_cloud_project on id=`).

## Examples

### List images of a cloud project

```sql
select
  id,
  name,
  type
from
  ovh_cloud_image
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```
