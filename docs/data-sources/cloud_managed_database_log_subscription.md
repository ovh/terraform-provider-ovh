---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_log_subscription (Data Source)

Use this data source to get a log subscription of a managed database service (mongodb, mysql, postgresql, valkey), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_database_log_subscription" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "postgresql"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  id           = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
}
```

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `engine` - (Required) Engine of the service: mongodb, mysql, postgresql, valkey.
* `service_id` - (Required) ID of the service.
* `id` - (Required) ID of the log subscription.

## Attributes Reference

The following attributes are exported:

* `created_at` - Creation date of the subscription.
* `kind` - Kind of logs forwarded.
* `ldp_service_name` - Name of the Logs Data Platform service of the stream.
* `resource` - Resource the logs come from.
  * `name` - Name of the resource.
  * `type` - Type of the resource.
* `stream_id` - ID of the Logs Data Platform stream receiving the logs.
* `updated_at` - Last update date of the subscription.
