# Table: ovh_cloud_loadbalancer

A load balancer distributes incoming traffic across a pool of backend servers.

The `ovh_cloud_loadbalancer` table can be used to query information about load balancers. You must specify which cloud project in the where clause (`where project_id=xxxx`). The region can optionally be specified in the where clause to filter results, or you can configure which regions to query using the `regions` parameter in your connection configuration.

**Note:** By default, this table will query all regions of the project offering the `octavialoadbalancer` service. To improve performance, you can:
- Specify a region in your query's where clause: `where region='GRA11'`
- Configure the `regions` parameter in your connection to limit which regions are queried globally

**Note:** This table covers the regional load balancers, exposed under `/cloud/project/{serviceName}/region/{regionName}/loadbalancing/loadbalancer`. These are the load balancers a managed Kubernetes cluster creates for a `Service` of type `LoadBalancer`. The older project wide load balancers, exposed under `/cloud/project/{serviceName}/loadbalancer`, are a different product and are not returned here.

## Examples

### List load balancers of a cloud project in all configured regions

```sql
select
  name,
  region,
  vip_address,
  operating_status,
  provisioning_status
from
  ovh_cloud_loadbalancer
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```

### List load balancers in a specific region

```sql
select
  name,
  vip_address,
  flavor_id,
  created_at
from
  ovh_cloud_loadbalancer
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and region='GRA11'
```

### List specific load balancer

```sql
select
  name,
  region,
  vip_address,
  vip_network_id,
  vip_subnet_id
from
  ovh_cloud_loadbalancer
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and region='GRA11'
  and id='287f3e89-70a9-4704-8c80-50661fa510a2'
```

### List load balancers that are not online

```sql
select
  name,
  region,
  operating_status,
  provisioning_status
from
  ovh_cloud_loadbalancer
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and operating_status != 'online'
```

### List load balancers exposed through a floating IP

```sql
select
  name,
  region,
  vip_address,
  floating_ip
from
  ovh_cloud_loadbalancer
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and floating_ip is not null
```
