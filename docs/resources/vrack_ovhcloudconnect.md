---
subcategory : "vRack"
---

# ovh_vrack_ovhcloudconnect

Attach an OVH Cloud Connect to the vrack.

## Example Usage

```terraform
resource "ovh_vrack_ovhcloudconnect" "vrack_ovhcloudconnect" {
  service_name      = "<vRack service name>"
  ovh_cloud_connect = "<OVH Cloud Connect service name>"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your vrack
* `ovh_cloud_connect` - (Required) Your OVH Cloud Connect service name.

## Attributes Reference

No additional attribute is exported.

## Import

Attachment of an OVH Cloud Connect and a vRack can be imported using the `service_name` (vRack identifier) and the `ovh_cloud_connect` (OVH Cloud Connect service name), separated by "/" E.g.,

```bash
$ terraform import ovh_vrack_ovhcloudconnect.myattach "<service_name>/<OVH Cloud Connect service name>"
```
