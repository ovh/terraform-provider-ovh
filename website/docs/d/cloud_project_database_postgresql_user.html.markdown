---
layout: "ovh"
page_title: "OVH: cloud_project_database_postgresql_user"
sidebar_current: "docs-ovh-datasource-cloud-project-database-postgresql-user"
description: |-
  Get information about a user of a postgresql cluster associated with a public cloud project.
---

# cloud_project_database_postgresql_user (Data Source)

Use this data source to get information about a user of a postgresql cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_postgresql_user" "pguser" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "pguser_roles" {
  value = data.ovh_cloud_project_database_postgresql_user.pguser.roles
}
```

## Argument Reference

* `service_name` - The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - Cluster ID

* `name` - Name of the user.

## Attributes Reference

The following attributes are exported:

* `created_at` - Date of the creation of the user.
* `id` - Public Cloud Database Service ID.
* `roles` - Roles the user belongs to.
* `status` - Current status of the user.
* `name` - Name of the user.
