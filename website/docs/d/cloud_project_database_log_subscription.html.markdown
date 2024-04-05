---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_log_subscription (Data Source)

Use this data source to get information about a log subscription for a cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_log_subscription" "subscription" {
    service_name = "VVV"
    engine       = "XXX"
    cluster_id   = "YYY"
    id           = "ZZZ"
}

output "subscription_ldp_name" {
  value = data.ovh_cloud_project_database_log_subscription.subscription.ldp_service_name
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `engine` - (Required) The database engine for which you want to retrieve a subscription. To get a full list of available engine visit.
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).
* `cluster_id` - (Required) Cluster ID.
* `id` - (Required) Id of the log subscription.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Creation date of the subscription.
* `engine` - See Argument Reference above.
* `id` - ID of the log subscription.
* `kind` - Log kind name of this subscription.
* `ldp_service_name` - Name of the destination log service.
* `resource_name` - Name of subscribed resource, where the logs come from.
* `resource_type` - Type of subscribed resource, where the logs come from.
* `service_name` - See Argument Reference above.
* `stream_id` - See Argument Reference above.
* `updated_at` - Last update date of the subscription.