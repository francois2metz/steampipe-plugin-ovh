---
organization: francois2metz
category: ["public cloud"]
brand_color: "#000e9c"
display_name: "OVH"
short_name: "ovh"
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

    # Go to https://www.ovh.com/auth/api/createToken to create your application key,
    # secret and the consumer key
    # For the rights, GET with the path *
    # application_key = "CitIbyantOosuzFu"
    # application_secret = "phoagDakOywytMibfetJidloidvuenVo"
    # consumer_key = "einbycsAnmachCeOkvabicdifAdofdon"

    # OVH Endpoint
    # 'ovh-eu' for OVH Europe API
    # 'ovh-us' for OVH US API
    # 'ovh-ca' for OVH Canada API
    # 'soyoustart-eu' for So you Start Europe API
    # 'soyoustart-ca' for So you Start Canada API
    # 'kimsufi-eu' for Kimsufi Europe API
    # 'kimsufi-ca' for Kimsufi Canada API
    endpoint = "ovh-eu"

    # List of regions to query. Supports wildcards.
    # Defaults to all regions if not specified.
    # regions = ["GRA", "SBG", "BHS"]
}
```

- `application_key` - Your OVH application key. Can also be set with the `OVH_APPLICATION_KEY` environment variable.
- `application_secret` - Your OVH application secret. Can also be set with the `OVH_APPLICATION_SECRET` environment variable.
- `consumer_key` - Your OVH consumer key. Can also be set with the `OVH_CONSUMER_KEY` environment variable.
- `endpoint` - The OVH API endpoint to use. Can also be set with the `OVH_ENDPOINT` environment variable.
- `regions` (Optional) - A list of regions to query. Supports wildcards (e.g., `GRA*` matches `GRA`, `GRA9`, etc.). If not specified, all regions will be queried. This can significantly improve query performance when you only need data from specific regions.

### Regions

The `regions` configuration parameter allows you to limit queries to specific OVH regions, which can significantly improve performance by reducing the number of API calls.

**Examples:**

Query a single region:
```hcl
regions = ["GRA"]
```

Query multiple specific regions:
```hcl
regions = ["GRA", "SBG", "BHS"]
```

Query all regions in a location using wildcards:
```hcl
regions = ["GRA*"]  # Matches GRA, GRA9, etc.
```

Query multiple locations:
```hcl
regions = ["GRA", "SBG", "BHS"]
```

Query all regions (default behavior):
```hcl
regions = ["*"]
```

**Note:** When `regions` is not specified in the configuration, it defaults to `["*"]` which queries all available regions.

## Get Involved

* Open source: https://github.com/francois2metz/steampipe-plugin-ovh
