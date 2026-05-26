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
