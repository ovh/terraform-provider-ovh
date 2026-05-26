---
subcategory: "VPS"
page_title: "OVHcloud: ovh_vps_disk"
description: |-
  Manages disk settings for an OVHcloud VPS disk.
---

# ovh_vps_disk

Manages disk settings for an OVHcloud VPS disk.

## Example Usage

```hcl
resource "ovh_vps_disk" "disk" {
  service_name        = "vps-xxxxxxxx.vps.ovh.net"
  id                  = 12345
  low_freespace_alert = 5
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service_name of your VPS.
* `id` - (Required) The ID of the disk.
* `low_freespace_alert` - (Optional) Alert threshold for low free space in GB.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the disk.
* `bandwith_limit` - Bandwidth limit of the disk in MB/s.
* `raid` - RAID type used on the disk.
* `size` - Size of the disk in GB.
* `state` - State of the disk.
* `type` - Type of disk (hdd┃ssd).
