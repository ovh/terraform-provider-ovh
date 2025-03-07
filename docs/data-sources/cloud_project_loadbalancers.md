---
subcategory : "Load Balancer (Public Cloud / Octavia)"
---

# ovh_cloud_project_loadbalancers

List your public cloud loadbalancers.

## Example Usage

```terraform
data "ovh_cloud_project_loadbalancers" "lbs" {
  service_name = "XXXXXX"
  region_name  = "XXX"
}
output "lbs" {
  value = data.ovh_cloud_project_loadbalancers.lbs
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

- `region_name` - (Required) Region of the loadbalancers.

## Attributes Reference

The following attributes are exported:

- `service_name` - ID of the public cloud project
- `region_name` - Region of the loadbalancers
- `loadbalancers` - List of loadbalancer
  - `id`: ID of the loadbalancer
  - `flavor_id`: ID of the flavor
  - `name`: Name of the loadbalancer
  - `region` - Region of the loadbalancer
  - `created_at` - Date of creation of the loadbalancer
  - `updated_at` - Last update date of the loadbalancer
  - `operating_status`: Operating status of the loadbalancer
  - `provisioning_status`: Provisioning status of the loadbalancer
  - `vip_address`: IP address of the Virtual IP
  - `vip_network_id`: Openstack ID of the network for the Virtual IP
  - `vip_subnet_id`: ID of the subnet for the Virtual IP
  - `floating_ip`: Information about the floating IP
    - `id`: ID of the floating IP
    - `ip`: Value of the floating IP
