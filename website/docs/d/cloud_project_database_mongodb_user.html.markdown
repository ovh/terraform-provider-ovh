---
layout: "ovh"
page_title: "OVH: cloud_project_database_mongodb_user"
sidebar_current: "docs-ovh-datasource-cloud-project-database-mongodb-user"
description: |-
  Get information about a user of a mongodb cluster associated with a public cloud project.
---

# ovh_cloud_project_database_mongodb_user (Data Source)

Use this data source to get information about a user of a mongodb cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_mongodb_user" "mongouser" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ@admin"
}

output "mongouser_roles" {
  value = data.ovh_cloud_project_database_mongodb_user.mongouser.roles
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `name` - (Required) Name of the user with the authentication database in the format name@authDB, for example: johndoe@admin

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `name` - Name of the user.
* `roles` - Roles the user belongs to
* `service_name` - Current status of the user.
* `status` - Current status of the user.
* `name` - See Argument Reference above.
