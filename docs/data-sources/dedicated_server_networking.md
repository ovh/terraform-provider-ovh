---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_networking (Data Source)

Use this data source to retrieve the network interfaces aggregation (bonding) of a dedicated server, including each interface's LACP fallback MAC address.

## Example Usage

```terraform
data "ovh_dedicated_server_networking" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

output "public_lacp_fallback_mac" {
  value = one([
    for iface in data.ovh_dedicated_server_networking.server.interfaces :
    iface.aggregation_fallback if iface.type == "public"
  ])
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your dedicated server.

## Attributes Reference

The following attributes are exported:

* `description` - Operation description.
* `status` - Operation status.
* `interfaces` - The list of interface aggregations, ordered by `type`. Each entry exports:
  * `macs` - The list of MAC addresses of the physical interfaces in this aggregation.
  * `type` - The network type of the interface (`public` or `vrack`).
  * `aggregation_fallback` - MAC address of the LACP fallback interface (the address the aggregation falls back to when the bond degrades to a single link).
