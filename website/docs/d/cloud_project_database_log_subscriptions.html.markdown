---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_log_subscriptions (Data Source)

Use this data source to get the list of log subscrition for a cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_log_subscriptions" "subscriptions" {
    service_name = "XXX"
    engine       = "YYY"
    cluster_id   = "ZZZ"
}

output "subscription_ids" {
  value = data.ovh_cloud_project_database_log_subscriptions.subscriptions.subscription_ids
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `engine` - (Required) The database engine for which you want to retrieve a subscription. To get a full list of available engine visit.
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).
* `cluster_id` - (Required) Cluster ID.

## Attributes Reference

`id` is set to the md5 sum of the list of all subscription ids. In addition,
the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `engine` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `subscription_ids` - The list of log subscription ids of the cluster associated with the project.