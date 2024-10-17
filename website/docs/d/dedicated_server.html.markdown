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

* `availability_zone` - Dedicated AZ localisation
* `boot_id` - Boot id of the server
* `boot_script` - Boot script of the server
* `urn` - URN of the dedicated server instance
* `commercial_range` - Dedicated server commercial range
* `datacenter` - Dedicated datacenter localisation (bhs1,bhs2,...)
* `efi_bootloader_path` - Path of the EFI bootloader of the dedicated server
* `ip` - Dedicated server ip (IPv4)
* `ips` - Dedicated server ip blocks
* `link_speed` - Link speed of the server
* `monitoring` - Icmp monitoring state
* `name` - Dedicated server name
* `new_upgrade_system`
* `no_intervention` - Prevent datacenter intervention
* `display_name` - Dedicated server display name
* `os` - Operating system
* `power_state` - Power state of the server (poweroff, poweron)
* `professional_use` - Does this server have professional use option
* `rack` - Rack id of the server
* `region` - Dedicated region localisation
* `rescue_mail` - Rescue mail of the server
* `rescue_ssh_key` - Public SSH Key used in the rescue mode
* `reverse` - Dedicated server reverse
* `root_device` - Root device of the server
* `server_id` - Server id
* `state` - Error, hacked, hackedBlocked, ok
* `support_level` - Dedicated server support level (critical, fastpath, gs, pro)
* `vnis` - The list of Virtualnetworkinterface associated with this server
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
