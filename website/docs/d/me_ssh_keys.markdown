---
subcategory : "Account Management"
---

# ovh_me_ssh_keys (Data Source)

Use this data source to retrieve list of names of the account's SSH keys.

## Example Usage

```hcl
data "ovh_me_ssh_keys" "mykeys" {}
```

## Argument Reference

This datasource takes no arguments.

## Attributes Reference

* `names` - The list of the names of all the SSH keys.
