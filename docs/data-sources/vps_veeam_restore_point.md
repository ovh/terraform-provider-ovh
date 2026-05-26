---
subcategory: "VPS"
page_title: "OVHcloud: ovh_vps_veeam_restore_point"
description: |-
  Get information about a specific Veeam restore point for an OVHcloud VPS.
---

# ovh_vps_veeam_restore_point (Data Source)

Use this data source to retrieve information about a specific Veeam restore point for an OVHcloud VPS.

## Example Usage

```hcl
data "ovh_vps_veeam_restore_point" "rp" {
  service_name = "vps-xxxxxxxx.vps.ovh.net"
  id           = 12345
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your VPS.
* `id` - (Required) The ID of the restore point.

## Attributes Reference

* `id` - The ID of the restore point.
* `creation_time` - Creation time of the restore point.
