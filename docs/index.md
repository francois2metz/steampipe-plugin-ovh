---
organization: francois2metz
category: ["public cloud"]
brand_color: "#000e9c"
display_name: "OVH"
short_name: "OVH"
description: "Steampipe plugin for querying OVH."
og_description: "Query OVH with SQL! Open source CLI. No DB required."
icon_url: "/images/plugins/francois2metz/ovh.svg"
og_image: "/images/plugins/francois2metz/ovh-social-graphic.png"
---

# OVH + Steampipe

[OVH](https://www.ovhcloud.com/) is a cloud computing company.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  id,
  description
from
  ovh_cloud_project
```

```
+----------------------------------+-------------+
| id                               | description |
+----------------------------------+-------------+
| 27c5a6d3dfez87893jfd88fdsfmvnqb8 | CloudWatt   |
| 81c5a6d3dfez87893jfd32fdsmpawq3z | Numergy     |
+----------------------------------+-------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/francois2metz/ovh/tables)**

## Get started

### Install

Download and install the latest OVH plugin:

```bash
steampipe plugin install francois2metz/ovh
```

### Configuration

Installing the latest ovh plugin will create a config file (`~/.steampipe/config/ovh.spc`) with a single connection named `ovh`:

```hcl
connection "ovh" {
    plugin = "francois2metz/ovh"

    application_key = ""
    application_secret = ""
    consumer_key = ""

    # Endpoint
    endpoint = "ovh-eu"
}
```

## Get Involved

* Open source: https://github.com/francois2metz/steampipe-plugin-ovh
