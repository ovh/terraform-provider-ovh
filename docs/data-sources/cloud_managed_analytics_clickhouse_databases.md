---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_clickhouse_databases (Data Source)

Use this data source to list the databases of a managed analytics service (clickhouse), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_analytics_clickhouse_databases" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
}
```

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `service_id` - (Required) ID of the service.

## Attributes Reference

The following attributes are exported:

* `engine` - Engine of the service: clickhouse.
* `databases` - The databases of the service.
  * `checksum` - Hash of the database, computed by the API.
  * `created_at` - Creation date of the database.
  * `default` - Whether the database is the default one, created with the service.
  * `id` - ID of the database.
  * `name` - Name of the database.
  * `updated_at` - Last update date of the database.
