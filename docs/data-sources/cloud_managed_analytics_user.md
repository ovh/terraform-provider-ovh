---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_user (Data Source)

Use this data source to get information about a user of a database cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_managed_analytics_user" "user" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
  name          = "UUU"
}

output "user_name" {
  value = data.ovh_cloud_managed_analytics_user.user.name
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required) The engine of the database cluster you want user information. To get a full list of available engine visit : [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases). Available engines:
  * `cassandra`
  * `kafka`
  * `kafkaConnect`
  * `mysql`
  * `grafana`

* `cluster_id` - (Required) Cluster ID

* `name` - (Required) Name of the user.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `engine` - See Argument Reference above.
* `id` - ID of the user.
* `service_name` - See Argument Reference above.
* `status` - Current status of the user.
* `name` - Name of the user.
