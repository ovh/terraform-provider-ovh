---
subcategory : "Account Management (IAM)"
---

# ovh_me_identity_group

Creates an identity group.

## Example Usage

```terraform
resource "ovh_me_identity_group" "my_group" {
  description = "Some custom description"
  name        = "my_group_name"
  role        = "NONE"
}
```

## Argument Reference

* `name` - Group name.
* `description` - Group description.
* `role` - Role associated with the group. Valid roles are ADMIN, REGULAR, UNPRIVILEGED, and NONE.

## Attributes Reference

* `urn` - URN of the user group, used when writing IAM policies
* `default_group` - Is the group a default and immutable one.
* `creation` - Creation date of this group.
* `last_update` - Date of the last update of this group.

## Import

Identity groups can be imported using their `name`:

```bash
$ terraform import ovh_me_identity_group.my_identity_group name
```