---
subcategory : "vRack"
---

# ovh_vrack_vps (Data Source)

Use this data source to retrieve information about a VPS attached to a vRack.

## Example Usage

```terraform
data "ovh_vrack_vps" "vrackvps" {
  service_name     = "pn-xxxxxxx"
  vps_service_name = "vpsXXXXX.vps.ovh.net"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service name of the vRack.
* `vps_service_name` - (Required) The service name of the VPS attached to the vRack.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `vps_service_name` - See Argument Reference above.
