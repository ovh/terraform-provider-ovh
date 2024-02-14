---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server (Data Source)

Use this data source to retrieve information about a dedicated server associated with your OVHcloud Account.

## Example Usage

```hcl
data "ovh_dedicated_server" "server" {
  service_name = "XXXXXX"
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your dedicated server.

## Attributes Reference

`id` is set with the service_name of the dedicated server.
In addition, the following attributes are exported:

* `boot_id` - boot id of the server
* `boot_script` - boot script of the server
* `urn` - URN of the dedicated server instance
* `commercial_range` - dedicated server commercial range
* `datacenter` - dedicated datacenter localisation (bhs1,bhs2,...)
* `ip` - dedicated server ip (IPv4)
* `ips` - dedicated server ip blocks
* `link_speed` - link speed of the server
* `monitoring` - Icmp monitoring state
* `name` - dedicated server name
* `os` - Operating system
* `professional_use` - Does this server have professional use option
* `rack` - rack id of the server
* `rescue_mail` - rescue mail of the server
* `reverse` - dedicated server reverse
* `root_device` - root device of the server
* `server_id` - your server id
* `state` - error, hacked, hackedBlocked, ok
* `support_level` - Dedicated server support level (critical, fastpath, gs, pro)
* `vnis` - the list of Virtualnetworkinterface assiociated with this server
  * `enabled` - VirtualNetworkInterface activation state
  * `mode` - VirtualNetworkInterface mode (public,vrack,vrack_aggregation)
  * `name` - User defined VirtualNetworkInterface name
  * `server_name` - Server bound to this VirtualNetworkInterface
  * `uuid` - VirtualNetworkInterface unique id
  * `vrack` - vRack name
  * `nics` - NetworkInterfaceControllers bound to this VirtualNetworkInterface
* `enabled_vrack_vnis` - List of enabled vrack VNI uuids
* `enabled_vrack_aggregation_vnis` - List of enabled vrack_aggregation VNI uuids
* `enabled_public_vnis` - List of enabled public VNI uuids
