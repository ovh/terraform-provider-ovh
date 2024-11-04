---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_redis_user (Data Source)

Use this data source to get information about a user of a redis cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_redis_user" "redis_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "redis_user_commands" {
  value = data.ovh_cloud_project_database_redis_user.redis_user.commands
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `name` - (Required) Name of the user

## Attributes Reference

The following attributes are exported:

* `categories` - Categories of the user.
* `channels` - Channels of the user.
* `cluster_id` - See Argument Reference above.
* `commands` - Commands of the user.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `keys` - Keys of the user.
* `name` - See Argument Reference above.
* `service_name` - Current status of the user.
* `status` - Current status of the user.
