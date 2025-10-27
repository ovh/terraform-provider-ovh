---
subcategory : "vRack"
---

# ovh_vracks (Data Source)

Use this data source to get the list of Vrack IDs available for your OVHcloud account.

## Example Usage

```terraform
data "ovh_vracks" "vracks" {}
```

## Argument Reference

This datasource takes no argument.

## Attributes Reference

The following attributes are exported:

* `urn` - The URN of the vps
* `result` - The list of vrack service name available for your OVHcloud account.
