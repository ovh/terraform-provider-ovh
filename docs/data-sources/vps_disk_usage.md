---
subcategory: "VPS"
page_title: "OVHcloud: ovh_vps_disk_usage"
description: |-
  Get disk usage statistics for an OVHcloud VPS disk.
---

# ovh_vps_disk_usage (Data Source)

Use this data source to retrieve disk usage statistics for an OVHcloud VPS disk.

## Example Usage

```hcl
data "ovh_vps_disk_usage" "usage" {
  service_name = "vps-xxxxxxxx.vps.ovh.net"
  id           = 12345
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your VPS.
* `id` - (Required) The ID of the disk.

## Attributes Reference

* `id` - The ID of the disk usage data source.
* `free` - Free space on the disk in GB.
* `used` - Used space on the disk in GB.
