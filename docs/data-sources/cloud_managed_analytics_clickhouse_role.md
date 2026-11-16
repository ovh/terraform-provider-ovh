---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_clickhouse_role (Data Source)

Use this data source to get a role of a managed analytics service (clickhouse), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_analytics_clickhouse_role" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  id           = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
}
```

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `service_id` - (Required) ID of the service.
* `id` - (Required) ID of the role.

## Attributes Reference

The following attributes are exported:

* `engine` - Engine of the service: clickhouse.
* `name` - Name of the role.
* `checksum` - Hash of the role, computed by the API.
* `created_at` - Creation date of the role.
* `updated_at` - Last update date of the role.
