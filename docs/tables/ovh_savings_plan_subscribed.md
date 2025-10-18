# Table: ovh_savings_plan_subscribed

List OVH Cloud Savings Plans subscribed for each service/project. Savings plans allow you to commit to a certain amount of compute resources for a period of time in exchange for discounted pricing.

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
  flavour,
  duration,
  start_date,
  end_date
from
  ovh_savings_plan_subscribed
where
  project_id = 'your-project-id-here';
```

### Get active savings plans with auto-renewal enabled

```sql
select
  savings_plan_id,
  display_name,
  size,
  flavour,
  duration,
  auto_renewal,
  period_end_action,
  end_date
from
  ovh_savings_plan_subscribed
where
  project_id = 'your-project-id-here'
  and status = 'active'
  and auto_renewal = true;
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
  and end_date < now() + interval '30 days'
  and status = 'active';
```

### Get a specific savings plan by ID

```sql
select
  *
from
  ovh_savings_plan_subscribed
where
  project_id = 'your-project-id-here'
  and savings_plan_id = 'sp_123456';
```