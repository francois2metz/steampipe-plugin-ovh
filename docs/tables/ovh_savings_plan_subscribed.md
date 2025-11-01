# Table: ovh_savings_plan_subscribed

List OVH Cloud Savings Plans subscribed for each Public Cloud project. Savings plans allow you to commit to a certain amount of compute resources for a period of time in exchange for discounted pricing.

The `ovh_savings_plan_subscribed` table can be used to query information about active savings plans for OVH Public Cloud projects.

**Important Notes:**
- You must specify a `project_id` in the WHERE clause to query this table
- The `project_id` corresponds to an OVH Public Cloud project ID (the UUID from `ovh_cloud_project`)

## Examples

### List all savings plans for a specific project

```sql
select
  savings_plan_id,
  display_name,
  status,
  size,
  flavor,
  period,
  start_date,
  end_date,
  termination_date
from
  ovh_savings_plan_subscribed
where
  project_id = 'your-project-id-here';
```

### Get savings plans with specific end action

```sql
select
  savings_plan_id,
  display_name,
  size,
  flavor,
  period,
  period_end_action,
  end_date
from
  ovh_savings_plan_subscribed
where
  project_id = 'your-project-id-here'
  and period_end_action = 'TERMINATE';
```

### Get savings plans ending soon

```sql
select
  savings_plan_id,
  display_name,
  status,
  end_date,
  period_end_action
from
  ovh_savings_plan_subscribed
where
  project_id = 'your-project-id-here'
  and end_date::date < current_date + interval '30 days'
  and status = 'ACTIVE';
```

### Analyze planned changes for savings plans

```sql
select
  display_name,
  status,
  planned_changes,
  jsonb_array_elements(planned_changes::jsonb) as planned_change
from
  ovh_savings_plan_subscribed
where
  project_id = 'your-project-id-here';
```

### Get all savings plans across multiple projects

```sql
WITH project_list AS (
  SELECT id, description FROM ovh_cloud_project
)
SELECT 
  p.description as project_name,
  p.id as project_id,
  sp.service_id,
  sp.savings_plan_id,
  sp.display_name,
  sp.status,
  sp.size,
  sp.flavor,
  sp.period,
  sp.offer_id,
  sp.period_end_action,
  sp.start_date,
  sp.end_date,
  sp.period_start_date,
  sp.period_end_date,
  sp.termination_date,
  sp.planned_changes
FROM project_list p
LEFT JOIN ovh_savings_plan_subscribed sp ON sp.project_id = p.id
WHERE sp.savings_plan_id IS NOT NULL  -- Only show projects with actual savings plans
ORDER BY p.description, sp.display_name;
```

### Get a specific savings plan by ID

```sql
select
  *
from
  ovh_savings_plan_subscribed
where
  project_id = 'your-project-id-here'
  and savings_plan_id = 'your-saving-plan-id-here';
```

## Schema

| Name | Type | Description |
|------|------|-------------|
| project_id | `string` | OVH Public Cloud project ID. |
| service_id | `int` | OVH service ID (internal billing ID). |
| savings_plan_id | `string` | Savings plan unique ID. |
| display_name | `string` | Human-readable plan name. |
| status | `string` | Plan status (active, terminated, etc.). |
| size | `int` | Number of resources covered by plan. |
| flavor | `string` | Savings Plan flavor (resource type). |
| period | `string` | Periodicity of the Savings Plan (duration, e.g., P1Y). |
| offer_id | `string` | Savings Plan commercial offer identifier. |
| period_end_action | `string` | Action performed when reaching the end of the period (REACTIVATE or TERMINATE). |
| start_date | `string` | Start date of the Savings Plan. |
| end_date | `string` | End date of the Savings Plan. |
| period_start_date | `string` | Start date of the current period. |
| period_end_date | `string` | End date of the current period. |
| termination_date | `string` | Date at which the Savings Plan is scheduled to be terminated (null if not scheduled for termination). |
| planned_changes | `json` | Changes planned on the Savings Plan. |

## Notes

- This table requires a `project_id` to be specified in the WHERE clause.
- The OVH API returns different response formats depending on whether savings plans exist for a project.
- Planned changes include status transitions (e.g., PENDING → ACTIVE → TERMINATED).
- All dates are returned as strings in ISO format.
- The `termination_date` field is null unless the plan is specifically scheduled for early termination.