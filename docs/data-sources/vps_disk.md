---
subcategory: "VPS"
page_title: "OVHcloud: ovh_vps_disk"
description: |-
  Get information about a specific disk attached to an OVHcloud VPS.
---

# ovh_vps_disk (Data Source)

Use this data source to retrieve information about a specific disk attached to an OVHcloud VPS.

## Example Usage

```hcl
data "ovh_vps_disk" "disk" {
  service_name = "vps-xxxxxxxx.vps.ovh.net"
  id           = 12345
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your VPS.
* `id` - (Required) The ID of the disk to retrieve.

## Attributes Reference

* `id` - The ID of the disk.
* `bandwith_limit` - Bandwidth limit of the disk in MB/s.
* `low_freespace_alert` - Alert threshold for low free space in GB.
* `raid` - RAID type used on the disk.
* `size` - Size of the disk in GB.
* `state` - State of the disk.
* `type` - Type of disk (hdd┃ssd).
