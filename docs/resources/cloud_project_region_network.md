---
subcategory : "Public Cloud Network"
---

# ovh_cloud_project_region_network

Creates a network in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_project_region_network" "net" {
   service_name = "XXXXXX"
   region_name  = "EU-SOUTH-LZ-MAD-A"
   name         = "Madrid Network"
   subnet       = {
      cidr              = "10.0.0.0/24"
      enable_dhcp       = true
      enable_gateway_ip = false
      ip_version        = 4
   }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the public cloud project
* `region_name` - (Required) Network region
* `name` - (Required) Name of the network
* `subnet` - (Required) Parameters to create a subnet
  * `cidr` - (Required) Subnet range in CIDR notation
  * `enable_dhcp` - (Required) Enable DHCP for the subnet
  * `enable_gateway_ip` - (Required) Set a gateway ip for the subnet
  * `ip_version` - (Required) IP version
  * `allocation_pools` - List of IP pools allocated in subnet
    * `start` - First IP for the pool (eg: 192.168.1.12)
    * `end` - Last IP for the pool (eg: 192.168.1.24)
  * `dns_name_servers` - DNS nameservers
  * `gateway_ip` - Gateway IP
  * `name` - Subnet name
  * `use_default_public_dnsresolver` - Use default DNS
  * `host_routes` - Host routes
    * `destination` - Host route destination (eg: 192.168.1.0/24)
    * `next_hop` - Host route next hop (eg: 192.168.1.254)
* `vlan_id` - VLAN ID, between 1 and 4000

## Attributes Reference

The following attributes are exported:

* `id` - Network id
* `service_name` - The id of the public cloud project
* `region_name` - Network region
* `region` - Network region returned by the API
* `name` - Name of the network
* `subnet` - Parameters to create a subnet
  * `cidr` - Subnet range in CIDR notation
  * `enable_dhcp` - Enable DHCP for the subnet
  * `enable_gateway_ip` - Set a gateway ip for the subnet
  * `ip_version` - IP version
  * `allocation_pools` - List of IP pools allocated in subnet
    * `start` - First IP for the pool (eg: 192.168.1.12)
    * `end` - Last IP for the pool (eg: 192.168.1.24)
  * `dns_name_servers` - DNS nameservers
  * `gateway_ip` - Gateway IP
  * `name` - Subnet name
  * `use_default_public_dnsresolver` - Use default DNS
  * `host_routes` - Host routes
    * `destination` - Host route destination (eg: 192.168.1.0/24)
    * `next_hop` - Host route next hop (eg: 192.168.1.254)
* `vlan_id` - VLAN ID, between 1 and 4000
* `visibility` - Network visibility

## Import

A network in a public cloud project can be imported using the `service_name`, `region_name` and `id` attributes. Using the following configuration:

```terraform
import {
  id = "<service_name>/<region_name>/<id>"
  to = ovh_cloud_project_region_network.test
}
```

You can then run:

```bash
$ terraform plan -generate-config-out=network.tf
$ terraform apply
```

The file `network.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above. See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.
