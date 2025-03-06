---
subcategory : "Web Cloud Private SQL"
---

# ovh_hosting_privatedatabase_user (Data Source)

Use this data source to retrieve information about an hosting privatedatabase user.

## Example Usage

```terraform
data "ovh_hosting_privatedatabase_user" "user" {
  service_name  = "XXXXXX"
  user_name     = "XXXXXX"
}
```

## Argument Reference

* `service_name` - The internal name of your private database
* `user_name` - User name

## Attributes Reference

`id` is set to `service_name`/`user_name`. In addition, the following attributes are exported.

* `creation_date` - Creation date of the database
* `databases` - Users granted to this database
  * `database_name` - Database's name linked to this user
  * `grant_type` - Grant of this user for this database
