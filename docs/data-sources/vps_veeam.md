---
subcategory: "VPS"
page_title: "OVHcloud: ovh_vps_veeam"
description: |-
  Get information about the Veeam backup service for an OVHcloud VPS.
---

# ovh_vps_veeam (Data Source)

Use this data source to retrieve information about the Veeam backup service for an OVHcloud VPS.

## Example Usage

```hcl
data "ovh_vps_veeam" "veeam" {
  service_name = "vps-xxxxxxxx.vps.ovh.net"
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your VPS.

## Attributes Reference

* `id` - The ID of the Veeam service.
* `backup_used` - Amount of backup storage used in GB.
* `ip` - IP address of the Veeam backup server.
* `state` - State of the Veeam service.

## Compatibility

This data source wraps `GET /vps/{serviceName}/veeam`. Live cross-region probing on 2026-05-26 shows
the endpoint is present in the **EU** and **CA** API schemas (`eu.api.ovh.com`,
`ca.api.ovh.com`) but **NOT** in the **US** schema (`api.us.ovhcloud.com`).

On a US-region VPS the OVHcloud API returns
`404: Got an invalid (or empty) URL`. Use this data source on EU or CA accounts,
or wait for OVH to expose this endpoint on US.
