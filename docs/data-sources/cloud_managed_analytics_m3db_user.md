---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_m3db_user (Data Source)

Use this data source to get information about a user of a M3DB cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_managed_analytics_m3db_user" "m3db_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "m3db_user_group" {
  value = data.ovh_cloud_managed_analytics_m3db_user.m3db_user.group
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `name` - (Required, Forces new resource) Name of the user.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `group` - See Argument Reference above.
* `name` - See Argument Reference above.
* `service_name` - Current status of the user.
* `status` - Current status of the user.
