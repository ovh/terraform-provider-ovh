---
subcategory : "Private Network"
---

# ovh_cloud_network_private_vracks (Data Source)

Use this data source to list the private networks (vRack) of a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_network_private_vracks" "networks" {
  service_name = "<public cloud project ID>"
}

output "network_names" {
  value = [for n in data.ovh_cloud_network_private_vracks.networks.networks : n.name]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.

## Attributes Reference

The following attributes are exported:

* `networks` - List of private networks. Each element exports:
  * `id` - Network ID.
  * `name` - Network name.
  * `description` - Network description.
  * `location` - Location of the network:
    * `region` - Region of the network.
  * `checksum` - Computed hash representing the current target specification value.
  * `created_at` - Creation date of the network.
  * `updated_at` - Last update date of the network.
  * `resource_status` - Network readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
  * `current_state` - Current state of the network:
    * `name` - Network name.
    * `description` - Network description.
    * `location` - Location details:
      * `region` - Region code.
