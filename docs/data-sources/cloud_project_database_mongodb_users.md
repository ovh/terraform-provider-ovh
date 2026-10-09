---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_mongodb_users (Data Source)

Use this data source to get the list of users of a mongodb cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_database_mongodb_users" "all" {
  service_name = "XXX"
  cluster_id   = "YYY"
}

output "mongo_user_names" {
  value = [for u in data.ovh_cloud_project_database_mongodb_users.all.users : u.name]
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

`id` is set to the md5 sum of the list of all user ids. In addition, the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `users` - The list of users of the mongodb cluster associated with the project, with their details.
  * `id` - ID of the user.
  * `name` - Name of the user with the authentication database in the format name@authDB.
  * `created_at` - Date of the creation of the user.
  * `roles` - Roles the user belongs to.
  * `status` - Current status of the user.
