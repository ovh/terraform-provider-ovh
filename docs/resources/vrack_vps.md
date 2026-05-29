---
subcategory : "vRack"
---

# ovh_vrack_vps

Attach a VPS to a vRack.

A VPS can be attached to only one vRack at a time. If you try to attach a VPS that is already attached, the API will return an error.

## Example Usage

```terraform
resource "ovh_vrack_vps" "vrackvps" {
  service_name     = "pn-xxxxxxx"
  vps_service_name = "vpsXXXXX.vps.ovh.net"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service name of the vRack. If omitted, the `OVH_VRACK_SERVICE` environment variable is used.

* `vps_service_name` - (Required) The service name of the VPS to attach to the vRack.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `vps_service_name` - See Argument Reference above.

## Import

Attachment of a VPS to a vRack can be imported using the `service_name` and `vps_service_name`, separated by "/" E.g.,

```bash
$ terraform import ovh_vrack_vps.vrackvps service_name/vps_service_name
```
