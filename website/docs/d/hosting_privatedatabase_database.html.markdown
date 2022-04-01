---
layout: "ovh"
page_title: "OVH: ovh_hosting_privatedatabase_database"
sidebar_current: "docs-ovh-datasource-hosting-privatedatabase-database"
description: |-
  Get information & status of an hosting privatedatabase.
---

# ovh_hosting_privatedatabase_database (Data Source)

Use this data source to retrieve information about an hosting privatedatabase.

## Example Usage

```hcl
data "ovh_hosting_privatedatabase_database" "mydatabase" {
  service_name  = "XXXXXX"
  database_name = "XXXXXX"
}
```

## Argument Reference

* `service_name` - The internal name of your private database
* `database_name` - Database name

## Attributes Reference

`id` is set to `service_name`/`database_name`. In addition, the following attributes are exported.

* `backup_time` - Time of the next backup (every day)
* `creation_date` - Creation date of the database
* `quota_used` - Space used by the database (in MB)
* `users` - Users granted to this database
  * `user_name` - User's name granted on this database
  * `grant_type` - Grant of this user for this database
