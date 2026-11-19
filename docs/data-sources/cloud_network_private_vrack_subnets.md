---
subcategory : "Private Network"
---

# ovh_cloud_network_private_vrack_subnets (Data Source)

Use this data source to list the subnets of a private network (vRack) in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_network_private_vrack_subnets" "subnets" {
  service_name = "<public cloud project ID>"
  network_id   = "<network ID>"
}

output "subnet_cidrs" {
  value = [for s in data.ovh_cloud_network_private_vrack_subnets.subnets.subnets : s.cidr]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `network_id` - (Required) Network ID of the parent private network.

## Attributes Reference

The following attributes are exported:

* `subnets` - List of subnets. Each element exports:
  * `id` - Subnet ID.
  * `name` - Subnet name.
  * `cidr` - CIDR address range for the subnet.
  * `location` - Location details:
    * `region` - Region code.
  * `description` - Subnet description.
  * `dhcp_enabled` - Whether DHCP is enabled on the subnet.
  * `dns_nameservers` - DNS nameservers for the subnet.
  * `gateway_ip` - Default gateway IP address.
  * `allocation_pools` - Ranges inside the subnet that addresses are handed out from automatically. When omitted on an IPv4 subnet with no gatewayIp, the whole CIDR is used except its first host address (for 10.0.0.0/24: 10.0.0.2 to 10.0.0.254). Each pool has:
    * `start` - Start IP address of the pool.
    * `end` - End IP address of the pool.
  * `host_routes` - Static routes handed to the instances of the subnet through DHCP:
    * `destination` - CIDR reached through the next hop.
    * `next_hop` - IP address, inside the subnet, that traffic to the destination is sent to.
  * `checksum` - Computed hash representing the current target specification value.
  * `created_at` - Creation date of the subnet.
  * `updated_at` - Last update date of the subnet.
  * `resource_status` - Subnet readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
  * `current_state` - Current state of the subnet:
    * `name` - Subnet name.
    * `cidr` - CIDR address range.
    * `description` - Subnet description.
    * `dhcp_enabled` - Whether DHCP is enabled.
    * `dns_nameservers` - Configured DNS nameservers.
    * `gateway_ip` - Default gateway IP address.
    * `host_routes` - Static host routes:
      * `destination` - Destination CIDR.
      * `next_hop` - Next hop IP address.
    * `location` - Location details:
      * `region` - Region code.
