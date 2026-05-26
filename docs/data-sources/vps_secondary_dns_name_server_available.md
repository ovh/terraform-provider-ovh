---
subcategory : "VPS"
---

# ovh_vps_secondary_dns_name_server_available (Data Source)

Returns the OVHcloud secondary DNS nameserver available for slaving zones from the given VPS. Useful when you need to publish the OVH NS records at your registrar before creating a `ovh_vps_secondary_dns_domain` resource.

## Example Usage

```terraform
data "ovh_vps_secondary_dns_name_server_available" "ns" {
  service_name = "vpsXXXXX.ovh.net"
}

output "ovh_secondary_ns_hostname" {
  value = data.ovh_vps_secondary_dns_name_server_available.ns.hostname
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.

## Attributes Reference

* `hostname` - Hostname of the OVHcloud secondary nameserver.
* `ip` - IPv4 of the nameserver.
* `ipv6` - IPv6 of the nameserver (may be empty).
