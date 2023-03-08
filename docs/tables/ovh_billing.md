# Table: ovh_billing

Billing is what you pay :).

The `ovh_billing` table can be used to query information about your billing information.

## Examples

### List bills

```sql
select * from ovh_billing;
```

### Get a bill

```sql
select * from ovh_billing where id = 'FRxxxxxxxx';
```
