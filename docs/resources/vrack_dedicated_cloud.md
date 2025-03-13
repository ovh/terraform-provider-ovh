---
subcategory : "vRack"
---

# ovh_vrack_dedicated_cloud

Attach a Dedicated Cloud to the vrack.

## Example Usage

```terraform
resource "ovh_vrack_dedicated_cloud" "vrack-dedicatedCloud" {
  service_name      = "<vRack service name>"
  dedicated_cloud   = "<Dedicated Cloud service name>"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your vrack
* `dedicated_cloud` - (Required) Your Dedicated Cloud service name

## Attributes Reference

No additional attribute is exported.

## Import

Attachment of a Dedicated Cloud and a vRack can be imported using the `service_name` (vRack identifier) and the `dedicated_cloud` (Dedicated Cloud service name), separated by "/" E.g.,

```bash
$ terraform import ovh_vrack_dedicated_cloud.myattach "<vRack service name>/<Dedicated Cloud service name>"
```
