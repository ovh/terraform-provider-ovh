---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_clickhouse_user (Data Source)

Use this data source to get information about a user of a clickhouse cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_database_clickhouse_user" "ch_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "ch_user_roles" {
  value = data.ovh_cloud_project_database_clickhouse_user.ch_user.roles
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `name` - (Required) Name of the user.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `roles` - Roles the user belongs to.
* `service_name` - See Argument Reference above.
* `status` - Current status of the user.
* `name` - Name of the user.
