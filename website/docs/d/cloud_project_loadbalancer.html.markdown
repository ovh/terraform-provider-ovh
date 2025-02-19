---
subcategory : "Load Balancer (Public Cloud / Octavia)"
---

# ovh_cloud_project_loadbalancer

Get the details of a public cloud project loadbalancer.

## Example Usage

```hcl
data "ovh_cloud_project_loadbalancer" "lb" {
  service_name = "XXXXXX"
  region_name  = "XXX"
  id           = "XXX"
}
output "lb" {
  value = data.ovh_cloud_project_loadbalancer.lb
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

- `region_name` - (Required) Region of the loadbalancer.

- `id`:  (Required) ID of the loadbalancer

## Attributes Reference

The following attributes are exported:

- `service_name` - ID of the public cloud project
- `region_name` - Region of the loadbalancer
- `id`:  ID of the loadbalancer
- `flavor_id`:  ID of the flavor
- `name`:  Name of the loadbalancer
- `created_at` - Date of creation of the loadbalancer
- `updated_at` - Last update date of the loadbalancer
- `operating_status`:  Operating status of the loadbalancer
- `provisioning_status`:   Provisioning status of the loadbalancer
- `vip_address`:  IP address of the Virtual IP
- `vip_network_id`:  Openstack ID of the network for the Virtual IP
- `vip_subnet_id`:   ID of the subnet for the Virtual IP
- `floating_ip`: Information about the floating IP
  - `id`: ID of the floating IP
  - `ip`: Value of the floating IP