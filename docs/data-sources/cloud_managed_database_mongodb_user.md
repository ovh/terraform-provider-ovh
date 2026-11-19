---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_mongodb_user (Data Source)

Use this data source to get a user of a managed database service (mongodb), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_database_mongodb_user" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  id           = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
}
```

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `service_id` - (Required) ID of the service.
* `id` - (Required) ID of the user.

## Attributes Reference

The following attributes are exported:

* `engine` - Engine of the service: mongodb.
* `admin` - Whether the user is the administrator of the service, created with it, which can't be deleted.
* `checksum` - Hash of the user, computed by the API.
* `created_at` - Creation date of the user.
* `credentials_last_updated_at` - Date of the last change of the password.
* `roles` - Roles of the user.
* `status` - Readiness of the user: CREATING, UPDATING, READY, ERROR or DELETING.
* `username` - Login name of the user.
