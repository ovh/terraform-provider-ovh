---
subcategory : "Cloud Storage"
---

# ovh_cloud_storage_file_share_network

Creates a file storage share network in a public cloud project. A share network binds a Neutron network and subnet to a region so that file shares can be attached to it.

Share networks are immutable: every argument forces the resource to be recreated when changed, and there is no update operation.

## Example Usage

```terraform
resource "ovh_cloud_storage_file_share_network" "network" {
  service_name = "<public cloud project ID>"
  name         = "my-share-network"
  description  = "Share network for my NFS shares"
  network_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subnet_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  region       = "GRA1"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `name` - (Required) Share network name. **Changing this value recreates the resource.**
* `network_id` - (Required) ID of the network backing the share network. **Changing this value recreates the resource.**
* `subnet_id` - (Required) ID of the subnet backing the share network. **Changing this value recreates the resource.**
* `region` - (Required) Region where the share network will be created. **Changing this value recreates the resource.**
* `description` - (Optional) Share network description. When omitted, this value is computed by the API (which may return an empty value). **Changing this value recreates the resource.**

## Attributes Reference

The following attributes are exported:

* `id` - Share network ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the share network.
* `updated_at` - Last update date of the share network.
* `resource_status` - Share network readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the file storage share network:
  * `name` - Share network name.
  * `description` - Share network description.
  * `network_id` - ID of the network backing the share network.
  * `subnet_id` - ID of the subnet backing the share network.
  * `location` - Current location:
    * `region` - Region.
    * `availability_zone` - Availability zone.

## Import

A cloud storage file share network can be imported using the `service_name` and `share_network_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_storage_file_share_network.network
  id = "<service_name>/<share_network_id>"
}
```

```bash
$ terraform import ovh_cloud_storage_file_share_network.network service_name/share_network_id
```
