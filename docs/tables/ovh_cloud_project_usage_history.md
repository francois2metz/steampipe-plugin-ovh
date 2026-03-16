# Table: ovh_cloud_project_usage_history

Retrieve historical usage and billing information for an OVHcloud Public Cloud project. This table provides access to past billing periods for trend analysis and cost tracking over time.

## Examples

### Basic historical usage information

```sql
select
  project_id,
  usage_id,
  usage_period_from,
  usage_period_to,
  grand_total_historical_price,
  comprehensive_total_historical_price
from
  ovh_cloud_project_usage_history
where project_id = '<your_project_id>';
```

### Historical usage by resource type

```sql
select
  project_id,
  usage_id,
  usage_period_from,
  total_instances_price,
  total_storage_price,
  total_volumes_price,
  total_snapshots_price,
  total_kubernetes_price
from
  ovh_cloud_project_usage_history
where project_id = '<your_project_id>'
order by
  usage_period_from desc;
```

### Detailed resource usage breakdown over time

```sql
select
  project_id,
  usage_id,
  usage_period_from,
  instances_count,
  volumes_count,
  storage_count,
  snapshots_count,
  total_instances_price,
  total_storage_price
from
  ovh_cloud_project_usage_history
where project_id = '<your_project_id>'
order by
  usage_period_from desc;
```

### Cost trend analysis

```sql
select
  usage_period_from,
  grand_total_historical_price,
  total_instances_price,
  total_storage_price,
  total_volumes_price,
  lag(grand_total_historical_price) over (order by usage_period_from) as previous_period_cost,
  grand_total_historical_price - lag(grand_total_historical_price) over (order by usage_period_from) as cost_change
from
  ovh_cloud_project_usage_history
where project_id = '<your_project_id>'
order by
  usage_period_from desc;
```

### Historical billing periods summary

```sql
select
  count(*) as total_periods,
  min(usage_period_from) as earliest_period,
  max(usage_period_from) as latest_period,
  sum(grand_total_historical_price) as total_historical_cost,
  avg(grand_total_historical_price) as average_period_cost
from
  ovh_cloud_project_usage_history
where project_id = '<your_project_id>';
```

## Schema

| Name | Type | Description |
|------|------|-------------|
| ai_hourly_usage | double | Hourly usage for AI and machine learning resources |
| ai_monthly_usage | double | Monthly usage for AI and machine learning resources |
| ai_total_price | double | Total price for AI and machine learning resources |
| archiving_hourly_usage | double | Hourly usage for archiving and cold storage |
| archiving_monthly_usage | double | Monthly usage for archiving and cold storage |
| archiving_total_price | double | Total price for archiving and cold storage |
| compute_hourly_usage | double | Hourly usage for compute resources (instances, CPUs) |
| compute_monthly_usage | double | Monthly usage for compute resources |
| compute_total_price | double | Total price for compute resources |
| database_hourly_usage | double | Hourly usage for database services |
| database_monthly_usage | double | Monthly usage for database services |
| database_total_price | double | Total price for database services |
| floating_ips_hourly_usage | double | Hourly usage for floating IP addresses |
| floating_ips_monthly_usage | double | Monthly usage for floating IP addresses |
| floating_ips_total_price | double | Total price for floating IP addresses |
| gateway_hourly_usage | double | Hourly usage for gateway services |
| gateway_monthly_usage | double | Monthly usage for gateway services |
| gateway_total_price | double | Total price for gateway services |
| licenses_hourly_usage | double | Hourly usage for software licenses |
| licenses_monthly_usage | double | Monthly usage for software licenses |
| licenses_total_price | double | Total price for software licenses |
| load_balancer_hourly_usage | double | Hourly usage for load balancer services |
| load_balancer_monthly_usage | double | Monthly usage for load balancer services |
| load_balancer_total_price | double | Total price for load balancer services |
| network_hourly_usage | double | Hourly usage for network resources (bandwidth, traffic) |
| network_monthly_usage | double | Monthly usage for network resources |
| network_total_price | double | Total price for network resources |
| notebooks_hourly_usage | double | Hourly usage for AI notebook services |
| notebooks_monthly_usage | double | Monthly usage for AI notebook services |
| notebooks_total_price | double | Total price for AI notebook services |
| object_storage_hourly_usage | double | Hourly usage for object storage |
| object_storage_monthly_usage | double | Monthly usage for object storage |
| object_storage_total_price | double | Total price for object storage |
| other_hourly_usage | double | Hourly usage for miscellaneous resources |
| other_monthly_usage | double | Monthly usage for miscellaneous resources |
| other_total_price | double | Total price for miscellaneous resources |
| period_end | timestamp | End date of the historical billing period |
| period_start | timestamp | Start date of the historical billing period |
| project_id | text | Unique identifier of the OVHcloud project |
| registry_hourly_usage | double | Hourly usage for container registry services |
| registry_monthly_usage | double | Monthly usage for container registry services |
| registry_total_price | double | Total price for container registry services |
| storage_hourly_usage | double | Hourly usage for storage resources (volumes, snapshots) |
| storage_monthly_usage | double | Monthly usage for storage resources |
| storage_total_price | double | Total price for storage resources |
| support_hourly_usage | double | Hourly usage for support services |
| support_monthly_usage | double | Monthly usage for support services |
| support_total_price | double | Total price for support services |
| total_hourly_cost | double | Total hourly cost across all resource types |
| total_monthly_cost | double | Total monthly cost across all resource types |
| total_price | double | Total price for the historical period |
| total_price_without_discount | double | Total price before any discounts applied |
| usage_id | text | Unique identifier for the specific usage period |

## Primary Key

- `project_id`
- `usage_id`