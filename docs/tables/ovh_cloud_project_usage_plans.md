# Table: ovh_cloud_project_usage_plans

Retrieve savings plan usage and cost information for an OVHcloud Public Cloud project. This table provides detailed information about savings plans, including costs, coverage, utilization, and subscription details.

## Examples

### Basic savings plan information

```sql
select
  project_id,
  period_from,
  period_to,
  total_savings,
  total_savings_currency,
  flavor,
  flavor_total_price
from
  ovh_cloud_project_usage_plans
where project_id = '<your_project_id>';
```

### Savings plan cost breakdown

```sql
select
  project_id,
  flavor,
  flat_fee_total_price,
  over_quota_quantity,
  over_quota_unit_price,
  flavor_saved_amount,
  flavor_total_price
from
  ovh_cloud_project_usage_plans
where project_id = '<your_project_id>';
```

### Usage coverage and utilization

```sql
select
  project_id,
  flavor,
  usage_period_coverage,
  usage_period_utilization,
  consumption_size,
  cumul_plan_size,
  subscription_size
from
  ovh_cloud_project_usage_plans
where project_id = '<your_project_id>';
```

### Subscription details

```sql
select
  project_id,
  subscription_id,
  plan_name,
  subscription_begin,
  subscription_end,
  subscription_size,
  flavor_saved_amount
from
  ovh_cloud_project_usage_plans
where project_id = '<your_project_id>';
```

## Schema

| Name | Type | Description |
|------|------|-------------|
| consumption_size | int | Number of instances consumed |
| cumul_plan_size | int | Cumulative plan size |
| flat_fee_currency | text | Currency for flat fee pricing |
| flat_fee_total_price | double | Total flat fee price for the flavor |
| flavor | text | Instance flavor type (e.g., b3-16) |
| flavor_saved_amount | double | Amount saved for this flavor |
| flavor_total_price | double | Total price for this flavor |
| flavors | jsonb | Complete flavors data (JSON array) |
| over_quota_quantity | int | Quantity of over-quota usage |
| over_quota_unit_price | double | Unit price for over-quota usage |
| period_from | timestamp | Start of the usage period |
| period_to | timestamp | End of the usage period |
| plan_name | text | Name of the savings plan |
| project_id | text | OVH Public Cloud project ID |
| subscription_begin | timestamp | Start date of the savings plan subscription |
| subscription_end | timestamp | End date of the savings plan subscription |
| subscription_id | text | ID of the savings plan subscription |
| subscription_size | int | Size of the savings plan subscription |
| total_savings | double | Total amount saved by using savings plans |
| total_savings_currency | text | Currency code for total savings |
| total_savings_text | text | Human-readable total savings amount |
| usage_period_coverage | text | Coverage percentage for the usage period |
| usage_period_utilization | text | Utilization percentage for the usage period |

## Primary Key

- `project_id`