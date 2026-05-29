---
subcategory : "VPS"
---

# ovh_vps_available_image (Data Source)

Use this data source to retrieve the details of a single OS image that is
available for installation on a given VPS.

## Example Usage

```terraform
data "ovh_vps_available_image" "img" {
  service_name = "vps-XXXXXX.vps.ovh.net"
  image_id     = "Debian 12"
}
```

## Argument Reference

* `service_name` - (Required) The service name of your VPS (e.g.
  `vps-XXXXXX.vps.ovh.net`).
* `image_id` - (Required) The ID of the image to look up. The full list
  of valid IDs for a given VPS can be obtained from
  `ovh_vps_available_images`.

## Attributes Reference

* `name` - The human-readable image name (e.g. `Debian 12`).
