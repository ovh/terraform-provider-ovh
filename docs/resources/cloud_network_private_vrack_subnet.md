---
subcategory : "Private Network"
---

# ovh_cloud_network_private_vrack_subnet

Creates a subnet in a private network (vRack) in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "xxxxxxxxxx"
  name         = "my-private-network"
  region       = "GRA1"
}

resource "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "my-subnet"
  cidr         = "10.0.0.0/24"
  region       = "GRA1"
  dhcp_enabled = true
  gateway_ip   = "10.0.0.1"

  dns_nameservers = [
    "213.186.33.99",
  ]

  allocation_pools {
    start = "10.0.0.2"
    end   = "10.0.0.254"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `network_id` - (Required) Network ID of the parent private network. **Changing this value recreates the resource.**
* `name` - (Required) Subnet name.
* `cidr` - (Required) CIDR address range for the subnet (e.g. `10.0.0.0/24`). **Changing this value recreates the resource.**
* `region` - (Required) Region where the subnet will be created. **Changing this value recreates the resource.**
* `description` - (Optional) Subnet description.
* `dhcp_enabled` - (Optional) Whether DHCP is enabled on the subnet.
* `dns_nameservers` - (Optional) List of DNS nameserver addresses.
* `gateway_ip` - (Optional) Default gateway IP address.
* `allocation_pools` - (Optional) IP address allocation pools:
  * `start` - (Required) Start IP address of the pool.
  * `end` - (Required) End IP address of the pool.

## Attributes Reference

The following attributes are exported:

* `id` - Subnet ID.
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

## Import

A cloud private network subnet can be imported using the `service_name`, `network_id` and `subnet_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_network_private_vrack_subnet.subnet
  id = "<service_name>/<network_id>/<subnet_id>"
}
```

```bash
$ terraform import ovh_cloud_network_private_vrack_subnet.subnet service_name/network_id/subnet_id
```
