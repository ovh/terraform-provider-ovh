---
subcategory : "VPS"
---

# ovh_vps_option (Data Source)

Use this data source to retrieve the state of a single option subscribed on a VPS.

This data source calls `GET /vps/{serviceName}/option/{option}`.

## Example Usage

```terraform
data "ovh_vps_option" "snapshot" {
  service_name = "vpsXXXXX.ovh.net"
  option       = "snapshot"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS (e.g. `vps-123456.ovh.net`).
* `option` - (Required) The option name. One of: `additionalDisk`, `automatedBackup`, `cpanel`, `ftpbackup`, `plesk`, `snapshot`, `veeam`, `windows`.

## Attributes Reference

The following attributes are exported:

* `state` - The subscription state of the option (e.g. `released`, `subscribed`).

~> Option lifecycle (subscribe / unsubscribe) is managed through the OVH billing/cart flow, not through this data source.
