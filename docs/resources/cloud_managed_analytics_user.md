---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_user

Creates a user in a managed analytics service (kafka, kafkaConnect), using the Cloud Databases API v2.

## Example Usage

```terraform
resource "ovh_cloud_managed_analytics_user" "app" {
  service_name   = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine         = "kafka"
  service_id     = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  name           = "app"
  password_reset = "2026-11-12T10:00:00Z"
}

output "app_password" {
  value     = ovh_cloud_managed_analytics_user.app.password
  sensitive = true
}
```

The API returns the password only at the creation and when `password_reset` resets it. An import leaves the password empty, and a `password_reset` in the configuration then resets it at the next apply. The administrator of the service, created with it, can't be deleted: destroying it only removes it from the state.

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used. Changing it recreates the user.
* `engine` - (Required) Engine of the service: kafka, kafkaConnect. Changing it recreates the user.
* `service_id` - (Required) ID of the service. Changing it recreates the user.
* `name` - (Required) Login name of the user. Changing it recreates the user.
* `password_reset` - (Optional) Arbitrary string to change to reset the password.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the user.
* `admin` - Whether the user is the administrator of the service, created with it, which can't be deleted.
* `checksum` - Hash of the user, computed by the API.
* `created_at` - Creation date of the user.
* `credentials_last_updated_at` - Date of the last change of the password.
* `password` - (Sensitive) Password of the user. The API returns it only at the creation and after a reset.
* `status` - Readiness of the user: CREATING, UPDATING, READY, ERROR or DELETING.
* `username` - Login name of the user.

## Timeouts

The creation, the update and the deletion wait up to 20 minutes for the user to be ready or deleted.

## Import

It can be imported using its `service_name`, `engine`, `service_id` and `id`, separated by "/".

```terraform
import {
  to = ovh_cloud_managed_analytics_user.app
  id = "<service_name>/<engine>/<service_id>/<id>"
}
```

```bash
$ terraform import ovh_cloud_managed_analytics_user.app service_name/engine/service_id/id
```
