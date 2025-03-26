---
subcategory : "Account Management (IAM)"
---

# ovh_iam_permissions_groups (Data Source)

Use this data source to retrieve all IAM permissions groups.

## Example Usage

```terraform
data "ovh_iam_permissions_groups" "groups" {}
```

## Attributes Reference

* `urns` - List of available permissions groups URNs.