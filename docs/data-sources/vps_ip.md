---
subcategory : "VPS"
---

# ovh_vps_ip (Data Source)

Use this data source to retrieve information about a single IP address
attached to a VPS.

## Example Usage

```terraform
data "ovh_vps_ip" "ip" {
  service_name = "vpsXXXXX.ovh.net"
  ip_address   = "192.0.2.1"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS (e.g. `vpsXXXXX.ovh.net`).
* `ip_address` - (Required) The IP address attached to the VPS.

## Attributes Reference

`id` is set to `service_name|ip_address`.

The following attributes are exported:

* `version` - IP version (`v4` or `v6`).
* `type` - IP type (`primary` or `additional`).
* `gateway` - Gateway of the IP, if any.
* `mac_address` - MAC address attached to the IP, if any.
* `geolocation` - Geolocation of the IP (read-only, set at order time).
* `reverse` - Current reverse DNS record for the IP.
