---
subcategory : "VPS"
---

# ovh_vps_options (Data Source)

Use this data source to list the options currently subscribed on a VPS, along with their state.

This data source calls `GET /vps/{serviceName}/option` to enumerate the subscribed option names, then fans out (capped at 4 concurrent workers) `GET /vps/{serviceName}/option/{option}` to gather each option's state.

## Example Usage

```terraform
data "ovh_vps_options" "opts" {
  service_name = "vpsXXXXX.ovh.net"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS (e.g. `vps-123456.ovh.net`).

## Attributes Reference

The following attributes are exported:

* `options` - The list of option names currently subscribed on the VPS. Possible values: `additionalDisk`, `automatedBackup`, `cpanel`, `ftpbackup`, `plesk`, `snapshot`, `veeam`, `windows`.
* `options_detail` - A list of objects providing per-option detail:
  * `name` - The option name.
  * `state` - The subscription state of the option (e.g. `released`, `subscribed`).

~> Option lifecycle (subscribe / unsubscribe) is managed through the OVH billing/cart flow, not through this data source.
