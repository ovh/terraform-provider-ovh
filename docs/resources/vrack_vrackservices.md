---
subcategory : "vRack"
---

# ovh_vrack_vrackservices

Attach a vrackServices to the vrack.

## Example Usage

```terraform
resource "ovh_vrack_vrackservices" "vrack_vrackservices" {
  service_name   = "<vRack service name>"
  vrack_services = "<vrackServices service name>"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your vrack
* `vrack_services` - (Required) Your vrackServices service name.

## Attributes Reference

No additional attribute is exported.

## Import

Attachment of a vrackServices and a vRack can be imported using the `service_name` (vRack identifier) and the `vrack_services` (vrackServices service name), separated by "/" E.g.,

```bash
$ terraform import ovh_vrack_vrackservices.myattach "<service_name>/<vrackServices service name>"
```
