---
subcategory : "Account Management"
---

# ovh_iam_reference_resource_type (Data Source)

Use this data source to list all the IAM resource types.

## Important
-> Using this resource requires that the account is enrolled in the OVHcloud [IAM beta](https://labs.ovhcloud.com/en/iam/) 

## Example Usage

```hcl
data "ovh_iam_reference_resource_type" "types" {
}
```

## Argument Reference

## Attributes Reference

* `id` - hash of the list of the resource types
* `types` - List of resource types
