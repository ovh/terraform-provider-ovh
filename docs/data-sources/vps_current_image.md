---
subcategory : "VPS"
---

# ovh_vps_current_image (Data Source)

Use this data source to retrieve the OS image currently installed on a
VPS. This is useful for drift detection or for keying other resources
off the live state of the VPS.

## Example Usage

```terraform
data "ovh_vps_current_image" "cur" {
  service_name = "vps-XXXXXX.vps.ovh.net"
}

output "current_image" {
  value = data.ovh_vps_current_image.cur.name
}
```

## Argument Reference

* `service_name` - (Required) The service name of your VPS (e.g.
  `vps-XXXXXX.vps.ovh.net`).

## Attributes Reference

* `id` - The ID of the image currently installed on the VPS.
* `name` - The human-readable name of the currently installed image
  (e.g. `Debian 12`).
