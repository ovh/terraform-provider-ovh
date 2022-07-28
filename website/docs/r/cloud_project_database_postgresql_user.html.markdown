---
layout: "ovh"
page_title: "OVH: cloud_project_database_postgresql_user"
sidebar_current: "docs-ovh-resource-cloud-project-database-postgresql-user"
description: |-
  Creates an user for a postgresql cluster associated with a public cloud project.
---

# cloud_project_database_postgresql_user

Creates an user for a postgresql cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database" "postgresql" {
  service_name  = "XXXX"
  engine        = "postgresql"
  cluster_id    = "ZZZZ"
}

resource "ovh_cloud_project_database_postgresql_user" "user" {
  service_name  = ovh_cloud_project_database.postgresql.service_name
  cluster_id    = ovh_cloud_project_database.postgresql.id
  name          = "johndoe"
  roles         = ["replication"]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - Cluster ID.

* `name` - Name of the user.

* `roles` - (Optional: if omit, default role) Roles the user belongs to. Possible values:
  * `["replication"]`
  * `[]` (default role)


## Attributes Reference

The following attributes are exported:

* `created_at` - Date of the creation of the user.
* `id` - Public Cloud Database Service ID.
* `password` - Password of the user.
* `roles` - Roles the user belongs to.
* `status` - Current status of the user.
* `name` - See Argument Reference above.

## Import

OVHcloud Managed postgresql clusters users can be imported using the `service_name`, `cluster_id` and `id` of the user, separated by "/" E.g.,

```
$ terraform import ovh_cloud_project_database_postgresql_user.my_user <service_name>/<cluster_id>/<id>