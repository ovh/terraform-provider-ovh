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

## Compatibility

This data source wraps `GET /vps/{serviceName}/veeam/restorePoints`. Live cross-region probing on 2026-05-26 shows
the endpoint is present in the **EU** and **CA** API schemas (`eu.api.ovh.com`,
`ca.api.ovh.com`) but **NOT** in the **US** schema (`api.us.ovhcloud.com`).

On a US-region VPS the OVHcloud API returns
`404: Got an invalid (or empty) URL`. Use this data source on EU or CA accounts,
or wait for OVH to expose this endpoint on US.
