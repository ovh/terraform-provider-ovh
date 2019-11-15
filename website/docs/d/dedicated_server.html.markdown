---
layout: "ovh"
page_title: "OVH: dedicated_server"
sidebar_current: "docs-ovh-datasource-dedicated-server-x"
description: |-
  Get information of a dedicated server associated with your OVH Account.
---

# ovh_dedicated_server

Use this data source to retrieve information about a dedicated server associated with 
your OVH Account.

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
* `commercial_range` - dedicater server commercial range
* `datacenter` - dedicated datacenter localisation (bhs1,bhs2,...)
* `ip` - dedicated server ip (IPv4)
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
  * `uuid` - VirtualNetworkInterface unique id
  * `vrack` - vRack name
  * `ncis` - NetworkInterfaceControllers bound to this VirtualNetworkInterface
* `enabled_vrack_vnis` - List of enabled vrack VNI uuids
* `enabled_vrack_aggregation_vnis` - List of enabled vrack_aggregation VNI uuids
* `enabled_public_vnis` - List of enabled public VNI uuids
