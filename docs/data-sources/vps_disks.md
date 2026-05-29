---
subcategory: "VPS"
page_title: "OVHcloud: ovh_vps_disks"
description: |-
  Get the list of disk IDs attached to an OVHcloud VPS.
---

# ovh_vps_disks (Data Source)

Use this data source to retrieve the list of disk IDs attached to an OVHcloud VPS.

## Example Usage

```hcl
data "ovh_vps_disks" "disks" {
  service_name = "vps-xxxxxxxx.vps.ovh.net"
}

output "disk_ids" {
  value = data.ovh_vps_disks.disks.result
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your VPS.

## Attributes Reference

* `id` - The ID of the data source.
* `result` - The list of disk IDs attached to this VPS.
