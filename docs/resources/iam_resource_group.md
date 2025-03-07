---
subcategory : "Account Management (IAM)"
---

# ovh_iam_resource_group (Resource)

Provides an OVHcloud IAM resource group.

## Example Usage

```terraform
resource "ovh_iam_resource_group" "my_resource_group" {
    name = "my_resource_group"
    resources = [
        "urn:v1:eu:resource:service1:service1-id",
        "urn:v1:eu:resource:service2:service2-id",
    ]
}
```

## Argument Reference

* `name`- Name of the resource group
* `resources`- Set of the URNs of the resources contained in the resource group. All urns must be ones of valid resources

## Attributes Reference

* `id`- Id of the resource group
* `owner`- Name of the account owning the resource group
* `created_at`- Date of the creation of the resource group
* `updated_at`- Date of the last modification of the resource group
* `read_only`- Marks that the resource group is not editable. Usually means that this is a default resource group created by OVHcloud
* `urn`- URN of the resource group, used when writing policies

## Import

Resource groups can be imported by using their id.

```bash
$ terraform import ovh_iam_resource_group.my_resource_group resource_group_id
```

-> Read only resource groups cannot be imported
