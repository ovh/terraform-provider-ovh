---
subcategory : "Private Network"
---

# ovh_cloud_network_private_vrack (Data Source)

Use this data source to retrieve information about a private network (vRack) in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_network_private_vrack" "network" {
  service_name = "<public cloud project ID>"
  id           = "<network ID>"
}

output "network_name" {
  value = data.ovh_cloud_network_private_vrack.network.name
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `id` - (Required) Network ID.

## Attributes Reference

The following attributes are exported:

* `name` - Network name.
* `location` - Location of the network:
  * `region` - Region of the network.
* `description` - Network description.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the network.
* `updated_at` - Last update date of the network.
* `resource_status` - Network readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the network:
  * `name` - Network name.
  * `description` - Network description.
  * `location` - Location details:
    * `region` - Region code.
