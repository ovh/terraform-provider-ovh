---
subcategory: "VPS"
page_title: "OVHcloud: ovh_vps_veeam_restore_points"
description: |-
  Get the list of Veeam restore point IDs for an OVHcloud VPS.
---

# ovh_vps_veeam_restore_points (Data Source)

Use this data source to retrieve the list of Veeam restore point IDs for an OVHcloud VPS.

## Example Usage

```hcl
data "ovh_vps_veeam_restore_points" "rps" {
  service_name = "vps-xxxxxxxx.vps.ovh.net"
}

output "restore_point_ids" {
  value = data.ovh_vps_veeam_restore_points.rps.result
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your VPS.

## Attributes Reference

* `id` - The ID of the data source.
* `result` - The list of Veeam restore point IDs.
