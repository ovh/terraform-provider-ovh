---
subcategory : "Account Management (IAM)"
---

# ovh_me_identity_group (Data Source)

Use this data source to retrieve information about an identity group.

## Example Usage

```terraform
data "ovh_me_identity_group" "my_group" {
  name = "my_group_name"
}
```

## Argument Reference

* `name` - Group name.

## Attributes Reference

* `urn` - Identity URN of the group.
* `description` - Group description.
* `role` - Role associated with the group. Valid roles are ADMIN, REGULAR, UNPRIVILEGED, and NONE.
* `default_group` - Is the group a default and immutable one.
* `creation` - Creation date of this group.
* `last_update` - Date of the last update of this group.
