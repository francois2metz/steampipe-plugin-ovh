# Table: ovh_cloud_image

Billing is what you pay :).

The `ovh_billing` table can be used to query information about your billing information.

## Examples

### List images of a cloud project

```sql
select * from ovh_billing;
```

Example:
```
+------------+---------------------------+----------------------------------------------------------------------------------------------------->
| id         | date                      | url                                                                                                 >
+------------+---------------------------+----------------------------------------------------------------------------------------------------->
| FRxxxxxxxx | 2017-04-19T08:24:35+01:00 | https://www.ovh.com/cgi-bin/order/bill.cgi?reference=FRxxxxxxxx&timestamp=1254896565&esign=a5f8e1d96>
| FRyyyyyyyy | 2022-06-02T11:47:40+02:00 | https://www.ovh.com/cgi-bin/order/bill.cgi?reference=FRyyyyyyyy&timestamp=4598989899&esign=58b48c9f5>
|
```

### Get on bill

```sql
select * from ovh_billing where id = 'FRxxxxxxxx';
```
