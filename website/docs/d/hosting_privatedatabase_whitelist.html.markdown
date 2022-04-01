---
layout: "ovh"
page_title: "OVH: ovh_hosting_privatedatabase_whitelist"
sidebar_current: "docs-ovh-datasource-hosting-privatedatabase-whitelist"
description: |-
  Get information & status of an hosting privatedatabase whitelist.
---

# ovh_hosting_privatedatabase_whitelist (Data Source)

Use this data source to retrieve information about an hosting privatedatabase whitelist.

## Example Usage

```hcl
data "ovh_hosting_privatedatabase_whitelist" "whitelist" {
  service_name  = "XXXXXX"
  ip            = "XXXXXX"
}
```

## Argument Reference

* `service_name` - The internal name of your private database
* `ip` - The whitelisted IP in your instance

## Attributes Reference

`id` is set to `service_name`/`ip`. In addition, the following attributes are exported.

* `creation_date` - Creation date of the database
* `last_update` - The last update date of this whitelist
* `name` - Custom name for your Whitelisted IP
* `service` - Authorize this IP to access service port
* `sftp` - Authorize this IP to access SFTP port
* `status` - Whitelist status