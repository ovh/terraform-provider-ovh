---
subcategory : "Private Network"
---

# ovh_cloud_network_private_vrack_subnet (Data Source)

Use this data source to retrieve information about a subnet in a private network (vRack) of a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = "<public cloud project ID>"
  network_id   = "<network ID>"
  id           = "<subnet ID>"
}

output "subnet_cidr" {
  value = data.ovh_cloud_network_private_vrack_subnet.subnet.cidr
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `network_id` - (Required) Network ID of the parent private network.
* `id` - (Required) Subnet ID.

## Attributes Reference

The following attributes are exported:

* `name` - Subnet name.
* `cidr` - CIDR address range for the subnet.
* `region` - Region of the subnet.
* `description` - Subnet description.
* `dhcp_enabled` - Whether DHCP is enabled on the subnet.
* `dns_nameservers` - DNS nameservers for the subnet.
* `gateway_ip` - Default gateway IP address.
* `allocation_pools` - IP address allocation pools:
  * `start` - Start IP address of the pool.
  * `end` - End IP address of the pool.
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
