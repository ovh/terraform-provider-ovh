---
subcategory : "Account Management (IAM)"
---

# ovh_me_identity_user

Creates an identity user.

~> **NOTE** The `groups` attribute manages group memberships inline. Alternatively, you can use
the [`ovh_me_identity_group_membership`](me_identity_group_membership) resource to manage memberships as standalone resources.
Do not use both approaches for the same user, as they will conflict with each other.

## Example Usage

```terraform
resource "ovh_me_identity_user" "my_user" {
  description = "Some custom description"
  email       = "my_login@example.com"
  group       = "DEFAULT"
  groups      = ["my_group", "another_group"]
  login       = "my_login"
  password    = "super-s3cr3t!password"
}
```

## Argument Reference

* `description` - User description.
* `email` - User's email.
* `group` - User's main group.
* `groups` - (Optional) Additional groups the user belongs to (other than the main group). Conflicts with `ovh_me_identity_group_membership`. Use one approach or the other, not both.
* `login` - User's login suffix.
* `password` - User's password.

## Attributes Reference

* `urn` - URN of the user, used when writing IAM policies
* `creation` - Creation date of this user.
* `last_update` - Last update of this user.
* `password_last_update` - When the user changed his password for the last time.
* `status` - Current user's status.

## Import

An identity user can be imported using the `login` E.g.,

```bash
$ terraform import ovh_me_identity_user.my_user login
```
