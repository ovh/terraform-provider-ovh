---
layout: "ovh"
page_title: "OVH: ovh_hosting_privatedatabase_user_grant"
sidebar_current: "docs-ovh-datasource-hosting-privatedatabase-user-grant"
description: |-
  Get information & status of an hosting privatedatabase user grant.
---

# ovh_hosting_privatedatabase_user_grant (Data Source)

Use this data source to retrieve information about an hosting privatedatabase user grant.

## Example Usage

```hcl
data "ovh_hosting_privatedatabase_user_grant" "user_grant" {
  service_name  = "XXXXXX"
  database_name = "XXXXXX"
  user_name     = "XXXXXX"
}
```

## Argument Reference

* `service_name` - The internal name of your private database
* `database_name` - The database name on which grant the user
* `user_name` - The user name

## Attributes Reference

`id` is set to `service_name`/`user_name`/`database_name`/`grant`. In addition, the following attributes are exported.

* `creation_date` - Creation date of the database
* `grant` - Grant name