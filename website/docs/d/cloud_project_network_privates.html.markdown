---
subcategory : "Public Cloud Network"
---

# ovh_cloud_project_network_privates

List public cloud project private networks.

## Example Usage

```hcl
data "ovh_cloud_project_network_privates" "private" {
  service_name = "XXXXXX"
}

output "private" {
  value = data.ovh_cloud_project_network_privates.private
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project.


## Attributes Reference

The following attributes are exported:

- `service_name` - ID of the public cloud project
- `networks` - List of network
  - `id` - ID of the network
  - `status` - Status of the network
  - `name` - Name of the network
  - `type` - Type of the network
  - `vlan_id` - VLAN ID of the network
  - `regions` - Information about the private network in the openstack region
    - `openstack_id` - Network ID on openstack region
    - `region` - Name of the region
    - `status` - Status of the network