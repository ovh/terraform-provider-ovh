---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_database

Creates a logical database in a managed database service (mysql, postgresql), using the Cloud Databases API v2.

## Example Usage

```terraform
resource "ovh_cloud_managed_database_database" "app" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "postgresql"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  name         = "app"
}
```

The default database of the service, created with it, can't be deleted: destroying it only removes it from the state.

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used. Changing it recreates the database.
* `engine` - (Required) Engine of the service: mysql, postgresql. Changing it recreates the database.
* `service_id` - (Required) ID of the service. Changing it recreates the database.
* `name` - (Required) Name of the database. Changing it recreates the database.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the database.
* `checksum` - Hash of the database, computed by the API.
* `created_at` - Creation date of the database.
* `default` - Whether the database is the default one, created with the service.
* `updated_at` - Last update date of the database.

## Import

It can be imported using its `service_name`, `engine`, `service_id` and `id`, separated by "/".

```terraform
import {
  to = ovh_cloud_managed_database_database.app
  id = "<service_name>/<engine>/<service_id>/<id>"
}
```

```bash
$ terraform import ovh_cloud_managed_database_database.app service_name/engine/service_id/id
```
