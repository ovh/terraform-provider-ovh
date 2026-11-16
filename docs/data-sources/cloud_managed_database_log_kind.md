---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_log_kind (Data Source)

Use this data source to get a log kind of a managed database service (mongodb, mysql, postgresql, valkey), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_database_log_kind" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "postgresql"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  name         = "customer_logs"
}
```

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `engine` - (Required) Engine of the service: mongodb, mysql, postgresql, valkey.
* `service_id` - (Required) ID of the service.
* `name` - (Required) Name of the log kind.

## Attributes Reference

The following attributes are exported:

* `additional_returned_fields` - Log fields the log kind adds.
* `created_at` - Creation date of the log kind.
* `display_name` - Display name of the log kind.
* `kind_id` - ID of the log kind.
* `updated_at` - Last update date of the log kind.
