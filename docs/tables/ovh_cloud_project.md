# Table: ovh_cloud_project

A cloud project is a way to regroup instance, storage, database, ... under a name.

The `ovh_cloud_project` table can be used to query information about cloud projects.

## Examples

### List projects

```sql
select
  id,
  description
from
  ovh_cloud_project
```
