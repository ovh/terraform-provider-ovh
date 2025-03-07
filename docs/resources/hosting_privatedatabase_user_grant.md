---
subcategory : "Web Cloud Private SQL"
---

# ovh_hosting_privatedatabase_user_grant

Add grant on a database in your private cloud database instance.

## Example Usage

```terraform
resource "ovh_hosting_privatedatabase_user_grant" "user_grant" {
  service_name  = "XXXXXX"
  user_name     = "terraform"
  database_name = "ovhcloud"
  grant         = "admin"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The internal name of your private database.
* `user_name` - (Required) User name used to connect on your databases.
* `database_name` - (Required) Database name where add grant.
* `grant` - (Required) Database name where add grant. Values can be:
  - admin
  - none
  - ro
  - rw

## Attributes Reference

The id is set to the value of `service_name`/`user_name`/`database_name`/`grant` .

## Import

OVHcloud database user's grant can be imported using the `service_name`, the `user_name`, the `database_name` and the `grant`, separated by "/" E.g.,

```
$ terraform import ovh_hosting_privatedatabase_user_grant.user service_name/user_name/database_name/grant
```
