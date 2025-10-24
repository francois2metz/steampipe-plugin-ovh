# Table: ovh_ceph

List of all the Ceph Clusters of your account.

The `ovh_ceph` table can be used to query information about your billing information.

## Examples


### List all clusters

```sql
SELECT
  id,
  region,
  size
FROM
 ovh_ceph
```

### List Active Clusters

```sql
SELECT
  id,
  region,
  size
FROM
 ovh_ceph
WHERE state LIKE '%ACTIVE%';
```

### Count clusters by versions

```sql
SELECT
  ceph_version,
  COUNT(ceph_version)
FROM
 ovh_ceph
GROUP BY ceph_version;
```

### Count clusters by status

```sql
SELECT
  status,
  COUNT(status)
FROM
 ovh_ceph
GROUP BY status;
```