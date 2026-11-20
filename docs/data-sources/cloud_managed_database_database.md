---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_database (Data Source)

Use this data source to get a database of a managed database service (mysql, postgresql), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_database_database" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "postgresql"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  id           = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
}
```

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `engine` - (Required) Engine of the service: mysql, postgresql.
* `service_id` - (Required) ID of the service.
* `id` - (Required) ID of the database.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the database.
* `checksum` - Hash of the database, computed by the API.
* `created_at` - Creation date of the database.
* `default` - Whether the database is the default one, created with the service.
* `updated_at` - Last update date of the database.
