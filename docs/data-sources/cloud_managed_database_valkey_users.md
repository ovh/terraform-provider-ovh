---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_valkey_users (Data Source)

Use this data source to list the users of a managed database service (valkey), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_database_valkey_users" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
}
```

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `service_id` - (Required) ID of the service.

## Attributes Reference

The following attributes are exported:

* `engine` - Engine of the service: valkey.
* `users` - The users of the service.
  * `admin` - Whether the user is the administrator of the service, created with it, which can't be deleted.
  * `categories` - Command categories the user is allowed or denied.
  * `channels` - Pub/Sub channels the user can access.
  * `checksum` - Hash of the user, computed by the API.
  * `commands` - Commands the user is allowed or denied.
  * `created_at` - Creation date of the user.
  * `credentials_last_updated_at` - Date of the last change of the password.
  * `id` - ID of the user.
  * `keys` - Key patterns the user can access.
  * `status` - Readiness of the user: CREATING, UPDATING, READY, ERROR or DELETING.
  * `username` - Login name of the user.
