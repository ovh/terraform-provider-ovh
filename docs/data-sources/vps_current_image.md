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

## Compatibility

This data source is part of the OVHcloud public VPS API schema, but the
underlying endpoint may not be implemented for every VPS lineup. Live
testing on a `vps-le-2-2-40` (2019v1) US instance returns
`404: Got an invalid (or empty) URL`. The endpoint is
exposed on EU and CA regions; availability on the US API depends on the VPS product generation.

If you receive a 404 from this data source, your VPS plan likely does
not expose this endpoint. See [the OVHcloud VPS API console](https://api.us.ovhcloud.com/console/#/vps)
for what's available on your specific service.
