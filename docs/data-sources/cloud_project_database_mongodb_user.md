---
subcategory : "Managed Databases"
---

~> **DEPRECATED:** Use `ovh_cloud_managed_database_mongodb_user` instead. This data source will be removed in the next major version.

# ovh_cloud_project_database_mongodb_user (Data Source)

Use this data source to get information about a user of a mongodb cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_database_mongodb_user" "mongo_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ@admin"
}

output "mongo_user_roles" {
  value = data.ovh_cloud_project_database_mongodb_user.mongo_user.roles
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `name` - (Required) Name of the user with the authentication database in the format name@authDB, for example: johndoe@admin

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `name` - See Argument Reference above.
* `roles` - Roles the user belongs to
* `service_name` - Current status of the user.
* `status` - Current status of the user.
