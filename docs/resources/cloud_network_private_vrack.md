---
subcategory : "Private Network"
---

# ovh_cloud_network_private_vrack

Creates a private network (vRack) in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "<public cloud project ID>"
  name         = "my-private-network"
  region       = "GRA1"
  description = "My private network"
  vlan_id     = 100
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `name` - (Required) Network name.
* `region` - (Required) Region where the network will be created. **Changing this value recreates the resource.**
* `description` - (Optional) Network description. **Changing this value recreates the resource.**
* `vlan_id` - (Optional) VLAN ID of the network (0-4096). Assigned by the API if not set. **Changing this value recreates the resource.** Not supported in localzone regions.

## Attributes Reference

The following attributes are exported:

* `id` - Network ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the network.
* `updated_at` - Last update date of the network.
* `resource_status` - Network readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the network:
  * `name` - Network name.
  * `description` - Network description.
  * `vlan_id` - VLAN ID.
  * `location` - Location details:
    * `region` - Region code.

## Import

A cloud private network can be imported using the `service_name` and `network_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_network_private_vrack.network
  id = "<service_name>/<network_id>"
}
```

```bash
$ terraform import ovh_cloud_network_private_vrack.network service_name/network_id
```
