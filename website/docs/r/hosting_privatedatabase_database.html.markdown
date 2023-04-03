---
subcategory : "Web Cloud Private SQL"
---

# ovh_hosting_privatedatabase_database

Create a new database on your private cloud database service.

## Example Usage

```hcl
resource "ovh_hosting_privatedatabase_database" "database" {
  service_name  = "XXXXXX"
  database_name = "XXXXXX"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The internal name of your private database.
* `database_name` - (Required) Name of your new database

## Attributes Reference

The id is set to the value of `service_name`/`database_name`.

## Import

OVHcloud Webhosting database can be imported using the `service_name` and the `database_name`, separated by "/" E.g.,

```
$ terraform import ovh_hosting_privatedatabase_database.database service_name/database_name
```