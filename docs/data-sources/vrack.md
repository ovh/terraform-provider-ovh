---
subcategory : "vRack"
---

# ovh_vrack (Data Source)

Use this data source to get Vrack information.

## Example Usage

```terraform
data "ovh_vrack" "my_vrack" {
    service_name = "pn-000000"
}
```

## Argument Reference

This datasource takes no argument.

## Attributes Reference

The following attributes are exported:

* `service_name` - The vrack service name for your OVHcloud account, e.g pn-000000.
* `name` - The name defined for the vrack
* `description` - The vrack description
- `iam` (Attributes) IAM resource metadata (see [below for nested schema](#nestedatt--iam))

<a id="nestedatt--iam"></a>
### Nested Schema for `iam`

Read-Only:

- `id` (String) Unique identifier of the resource
- `state` (String) Resource state
- `urn` (String) Unique resource name used in policies