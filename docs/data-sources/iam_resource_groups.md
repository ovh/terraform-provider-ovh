---
subcategory : "Account Management (IAM)"
---

# ovh_iam_resource_groups (Data Source)

Use this data source to list the existing IAM policies of an account.

## Example Usage

```terraform
data "ovh_iam_resource_groups" "my_groups" {
}
```

## Argument Reference

## Attributes Reference

* `id` - Hash of the list of the resource groups IDs.
* `resource_groups` - List of the resource groups IDs.
