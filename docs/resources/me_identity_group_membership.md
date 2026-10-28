---
subcategory : "Account Management (IAM)"
---

# ovh_me_identity_group_membership

Manages the membership of an identity user in an identity group as a standalone resource.

~> **NOTE** This resource is an alternative to the `groups` attribute on `ovh_me_identity_user`.
Do not use both `ovh_me_identity_group_membership` and the `groups` attribute on `ovh_me_identity_user` for the same user, as they will conflict with each other.

## Example Usage

```terraform
resource "ovh_me_identity_group" "my_group" {
  description = "My custom group"
  name        = "my_group"
  role        = "NONE"
}

resource "ovh_me_identity_user" "my_user" {
  description = "My custom user"
  email       = "my_user@example.com"
  group       = "DEFAULT"
  login       = "my_user"
  password    = "super-s3cr3t!password"
}

resource "ovh_me_identity_group_membership" "my_membership" {
  login = ovh_me_identity_user.my_user.login
  group = ovh_me_identity_group.my_group.name
}
```

## Argument Reference

* `login` - (Required, Forces new resource) The login of the identity user to add to the group.
* `group` - (Required, Forces new resource) The name of the identity group to add the user to.

## Import

An identity group membership can be imported using `login/group`:

```bash
$ terraform import ovh_me_identity_group_membership.example my_login/my_group
```
