---
layout: "ovh"
page_title: "OVH: ovh_hosting_privatedatabase_user"
sidebar_current: "docs-ovh-resource-hosting-privatedatabase-user"
description: |-
  Create a new user on your private cloud database instance.
---

# ovh_hosting_privatedatabase_user

Create a new user on your private cloud database instance.

## Example Usage

```hcl
resource "ovh_hosting_privatedatabase_user" "user" {
  service_name  = "XXXXXX"
  password      = "XXXXXX"
  user_name     = "XXXXXX"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The internal name of your private database.
* `password` - (Required) Password for the new user (alphanumeric, minimum one number and 8 characters minimum)
* `user_name` - (Required) User name used to connect on your databases

## Attributes Reference

The id is set to the value of `service_name`/`user_name`.

## Import

OVHcloud database user can be imported using the `service_name` and the `user_name`, separated by "/" E.g.,

```
$ terraform import ovh_hosting_privatedatabase_user.user service_name/user_name
```