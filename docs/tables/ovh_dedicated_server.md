# Table: ovh_dedicated_server

Lists all OVH dedicated servers with their hardware and configuration details, providing a comprehensive inventory of your dedicated server infrastructure.

## Examples

### Basic server inventory

```sql
SELECT
  name,
  server_id,
  ip,
  state,
  power_state,
  os,
  datacenter,
  commercial_range
FROM
  ovh_dedicated_server
ORDER BY
  datacenter, name;
```

### Servers by operating system

```sql
SELECT
  os,
  COUNT(*) as server_count,
  STRING_AGG(name, ', ') as servers
FROM
  ovh_dedicated_server
GROUP BY
  os
ORDER BY
  server_count DESC;
```

### Servers by datacenter and commercial range

```sql
SELECT
  datacenter,
  commercial_range,
  COUNT(*) as server_count,
  ARRAY_AGG(name) as server_names
FROM
  ovh_dedicated_server
GROUP BY
  datacenter, commercial_range
ORDER BY
  datacenter, server_count DESC;
```

### High-speed network servers

```sql
SELECT
  name,
  ip,
  link_speed,
  datacenter,
  commercial_range
FROM
  ovh_dedicated_server
WHERE
  link_speed >= 10000
ORDER BY
  link_speed DESC;
```

### Server power states monitoring

```sql
SELECT
  power_state,
  COUNT(*) as server_count,
  STRING_AGG(name, ', ') as servers
FROM
  ovh_dedicated_server
GROUP BY
  power_state;
```

### Servers with monitoring enabled

```sql
SELECT
  name,
  ip,
  datacenter,
  monitoring,
  support_level,
  professional_use
FROM
  ovh_dedicated_server
WHERE
  monitoring = true
ORDER BY
  support_level, datacenter;
```

### Professional vs personal use servers

```sql
SELECT
  professional_use,
  support_level,
  COUNT(*) as server_count,
  AVG(link_speed) as avg_link_speed
FROM
  ovh_dedicated_server
GROUP BY
  professional_use, support_level
ORDER BY
  professional_use DESC, support_level;
```

### Get specific server details

```sql
SELECT
  *
FROM
  ovh_dedicated_server
WHERE
  name = 'ns3013242.ip-57-128-124.eu';
```
