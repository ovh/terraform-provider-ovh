---
subcategory : "Account Management (IAM)"
---

# ovh_iam_resource_group (Data Source)

Use this data source get details about a resource group.

## Example Usage

```terraform
data "ovh_iam_resource_group" "my_resource_group" {
    id = "my_resource_group_id"
}
```

## Argument Reference

* `id`- Id of the resource group

## Attributes Reference

* `name`- Name of the resource group
* `resources`- Set of the URNs of the resources contained in the resource group
* `owner`- Name of the account owning the resource group
* `created_at`- Date of the creation of the resource group
* `updated_at`- Date of the last modification of the resource group
* `read_only`- Marks that the resource group is not editable. Usually means that this is a default resource group created by OVHcloud
* `urn`- URN of the resource group, used when writing policies
