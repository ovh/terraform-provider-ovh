---
subcategory : "Private Network"
---

# ovh_cloud_project_network_private_subnets

List public cloud project subnets of a private network.

## Example Usage

```terraform
data "ovh_cloud_project_network_private_subnets" "private" {
  service_name = "XXXXXX"
  network_id   = "XXXXXX"
}
output "private" {
  value = data.ovh_cloud_project_network_private_subnets.private
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project.

- `network_id`: (Required) ID of the network

## Attributes Reference

The following attributes are exported:

- `service_name` - ID of the public cloud project
- `network_id` - ID of the network
- `subnets` - List of subnets
  - `id` - ID of the subnet
  - `cidr` - CIDR of the subnet
  - `dhcp_enabled` - Whether or not if DHCP is enabled for the subnet
  - `gateway_ip` - Gateway IP of the subnet
  - `ip_pools` - List of ip pools allocated in the subnet
    - `dhcp` - Whether or not if DHCP is enabled
    - `start` - First IP for this region (eg: 192.168.1.12)
    - `end` - Last IP for this region (eg: 192.168.1.24)
    - `region` - Region associated to the subnet
    - `network` - Global network with cidr (eg: 192.168.1.0/24)
