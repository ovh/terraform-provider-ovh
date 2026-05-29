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

## Compatibility

This data source wraps `GET /vps/{serviceName}/veeam/restorePoints/{id}`. Live cross-region probing on 2026-05-26 shows
the endpoint is present in the **EU** and **CA** API schemas (`eu.api.ovh.com`,
`ca.api.ovh.com`) but **NOT** in the **US** schema (`api.us.ovhcloud.com`).

On a US-region VPS the OVHcloud API returns
`404: Got an invalid (or empty) URL`. Use this data source on EU or CA accounts,
or wait for OVH to expose this endpoint on US.
