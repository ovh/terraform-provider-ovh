---
subcategory : "Private Network"
---

# ovh_cloud_project_network_private

Get the details of a public cloud project private network.

## Example Usage

```terraform
data "ovh_cloud_project_network_private" "private" {
  service_name = "XXXXXX"
  network_id           = "XXX"
}
output "private" {
  value = data.ovh_cloud_project_network_private.private
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project.

- `network_id`: (Required) ID of the network

## Attributes Reference

The following attributes are exported:

- `service_name` - ID of the public cloud project
- `network_id`: ID of the network
- `status`: Status of the network
- `name`: Name of the network
- `type`: Type of the network
- `vlan_id`: VLAN ID of the network
- `regions`: Information about the private network in the openstack region
  - `openstack_id`: Network ID on openstack region
  - `region`: Name of the region
  - `status`: Status of the network
