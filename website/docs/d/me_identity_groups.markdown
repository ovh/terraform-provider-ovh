---
subcategory : "Account Management"
---

# ovh_me_identity_groups (Data Source)

Use this data source to retrieve the list of the account's identity groups

## Example Usage

```hcl
data "ovh_me_identity_groups" "groups" {}
```

## Argument Reference

This datasource takes no arguments.

## Attributes Reference

* `groups` - The list of the group names of all the identity groups.
