---
layout: "ovh"
page_title: "OVH: me_identity_user"
sidebar_current: "docs-ovh-datasource-identity-user-x"
description: |-
  Get information about identity User.
---

# ovh_me_identity_user

Use this data source to retrieve information about an identity user.

## Example Usage

```hcl
data "ovh_me_identity_user" "my_user" {
   user = "my_user_login"
}
```

## Argument Reference

* `user` - (Required) User's login.

## Attributes Reference

* `login` - User's login suffix.
* `creation` - Creation date of this user.
* `description` - User description.
* `email` - User's email.
* `group` - User's group.
* `last_update` - Last update of this user.
* `password_last_update` - When the user changed his password for the last time.
* `status` - Current user's status.
