---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_users (Data Source)

Use this data source to get the list of users of a database cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_managed_analytics_users" "users" {
  service_name  = "XXXX"
  engine        = "YYYY"
  cluster_id    = "ZZZ"
}

output "user_ids" {
  value = data.ovh_cloud_managed_analytics_users.users.user_ids
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required) The engine of the database cluster you want to list users. To get a full list of available engine visit: [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

`id` is set to the md5 sum of the list of all user ids. In addition, the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `engine` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `user_ids` - The list of users ids of the database cluster associated with the project.
