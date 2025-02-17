---
subcategory : "vRack"
---

# ovh_vrack_ipv6

Attach an IPv6 block to a VRack.


## Example Usage

```hcl
resource "ovh_vrack_ipv6" "vrack_block" {
  service_name = "<vRack service name>"
  block        = "<ipv6 block>"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your vrack
* `block` - (Required) Your IPv6 block.
    
## Attributes Reference

No additional attribute is exported.


## Import

Attachment of an IPv6 block and a VRack can be imported using the `service_name` (vRack identifier) and the `block` (IPv6 block), separated by "," E.g.,

```bash
$ terraform import ovh_vrack_ipv6.myattach "<service_name>,<block>"
```
