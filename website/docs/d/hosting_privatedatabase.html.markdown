---
layout: "ovh"
page_title: "OVH: hosting_privatedatabase"
sidebar_current: "docs-ovh-datasource-hosting-privatedatabase"
description: |-
  Get information & status of an hosting database.
---

# ovh_hosting_privatedatabase (Data Source)

Use this data source to retrieve information about an hosting database.

## Example Usage

```hcl
data "ovh_hosting_privatedatabase" "database" {
 service_name  = "XXXXXX"
}
```

## Argument Reference

* `service_name` - The service name

## Attributes Reference

`id` is set to database service_name. In addition, the following attributes are exported.

* `cpu` - Number of cpu on your private database
* `datacenter` - Datacenter where this private database is located
* `display_name` - Name displayed in customer panel for your private database
* `hostname` - Private database hostname
* `hostname_ftp` - Private database ftp hostname
* `infrastructure` - Infrastructure where service was stored
* `offer` - Type of the private database offer
* `port` - Private database service port
* `port_ftp` - Private database ftp port
* `quota_size` - Space allowed (in MB) on your private database
* `quota_used` - Sapce used (in MB) on your private database
* `ram` - Amount of ram (in MB) on your private database
* `server` - Private database server name
* `state` - Private database state
* `version` - Private database available versions
* `version_label` - Private database version label
* `version_number` - Private database version number
