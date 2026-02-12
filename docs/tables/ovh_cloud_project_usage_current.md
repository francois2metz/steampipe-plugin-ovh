# Table: ovh_cloud_project_usage_current

Retrieve current usage and billing information for an OVHcloud Public Cloud project. This table provides real-time consumption and cost data for the ongoing billing period.

## Examples

### Basic current usage information

```sql
select
  project_id,
  last_update,
  total_instances_price,
  total_storage_price,
  total_volumes_price,
  total_snapshots_price
from
  ovh_cloud_project_usage_current
where project_id = '<your_project_id>';
```

### Current usage by resource type

```sql
select
  project_id,
  last_update,
  instances_count,
  total_instances_price,
  storage_count,
  total_storage_price,
  volumes_count,
  total_volumes_price
from
  ovh_cloud_project_usage_current
where project_id = '<your_project_id>';
```

### Detailed usage breakdown with counts

```sql
select
  project_id,
  last_update,
  instances_count,
  volumes_count,
  storage_count,
  snapshots_count,
  total_instances_price,
  total_storage_price
from
  ovh_cloud_project_usage_current
where project_id = '<your_project_id>';
```

### Comprehensive cost analysis with savings

```sql
select
  project_id,
  period_start,
  total_price,
  compute_total_price,
  storage_total_price,
  network_total_price,
  ai_total_price,
  database_total_price,
  licenses_total_price,
  support_total_price,
  other_total_price,
  total_savings
from
  ovh_cloud_project_usage_current
order by
  total_price desc;
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
| period_end | timestamp | End date of the current billing period |
| period_start | timestamp | Start date of the current billing period |
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
| total_price | double | Total price for the current period |
| total_price_without_discount | double | Total price before any discounts applied |
| total_savings | double | Total savings applied to the current usage |
| volume_snapshots_hourly_usage | double | Hourly usage for volume snapshots |
| volume_snapshots_monthly_usage | double | Monthly usage for volume snapshots |
| volume_snapshots_total_price | double | Total price for volume snapshots |

## Primary Key

- `project_id`
- `period_start`