---
subcategory: "VPS"
page_title: "OVHcloud: ovh_vps_vnc"
description: |-
  Manages VNC access for an OVHcloud VPS.
---

# ovh_vps_vnc

Manages VNC access for an OVHcloud VPS.

## Example Usage

```hcl
resource "ovh_vps_vnc" "vnc" {
  service_name = "vps-xxxxxxxx.vps.ovh.net"
}

output "vnc_url" {
  value = "${ovh_vps_vnc.vnc.host}:${ovh_vps_vnc.vnc.port}"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service_name of your VPS.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VNC resource.
* `host` - Hostname of the VNC server.
* `password` - Password for VNC access.
* `port` - Port of the VNC server.
* `type` - Type of VNC connection (novnc┃vnc).

## Compatibility

This resource wraps `POST /vps/{serviceName}/openConsoleAccess`. Live cross-region probing on 2026-05-26 shows
the endpoint is present in the **EU** and **CA** API schemas (`eu.api.ovh.com`,
`ca.api.ovh.com`) but **NOT** in the **US** schema (`api.us.ovhcloud.com`).

On a US-region VPS the OVHcloud API returns
`404: Got an invalid (or empty) URL`. Use this resource on EU or CA accounts,
or wait for OVH to expose this endpoint on US.
