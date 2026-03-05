---
subcategory : "Cloud Storage"
---

# ovh_cloud_storage_block_volume

Creates a block storage volume in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "xxxxxxxxxx"
  name         = "my-volume"
  size         = 10
  region       = "GRA1"
  volume_type  = "CLASSIC"
}
```

### Create from backup

```terraform
resource "ovh_cloud_storage_block_volume" "restored" {
  service_name = "xxxxxxxxxx"
  name         = "my-restored-volume"
  size         = 10
  region       = "GRA1"
  volume_type  = "CLASSIC"

  create_from {
    backup_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `name` - (Required) Volume name.
* `size` - (Required) Size of the volume in GB.
* `region` - (Required) Region where the volume will be created. **Changing this value recreates the resource.**
* `volume_type` - (Required) Volume type (`CLASSIC`, `HIGH_SPEED`, `HIGH_SPEED_GEN2`). **Changing this value recreates the resource.**
* `create_from` - (Optional) Source to create the volume from. **Changing this value recreates the resource.**
  * `backup_id` - (Optional) Identifier of a backup to restore the volume from.

## Attributes Reference

The following attributes are exported:

* `id` - Volume ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the volume.
* `updated_at` - Last update date of the volume.
* `resource_status` - Volume readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the block storage volume:
  * `location` - Current location:
    * `region` - Region.
  * `name` - Volume name.
  * `size` - Size of the volume in GB.
  * `volume_type` - Volume type (`CLASSIC`, `HIGH_SPEED`, `HIGH_SPEED_GEN2`).
  * `status` - Volume status (`AVAILABLE`, `IN_USE`, `CREATING`, `DELETING`, `ATTACHING`, `DETACHING`, `EXTENDING`, `ERROR`, `ERROR_DELETING`, `ERROR_BACKING_UP`, `ERROR_RESTORING`, `ERROR_EXTENDING`).

## Import

A cloud storage block volume can be imported using the `service_name` and `volume_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_storage_block_volume.volume
  id = "<service_name>/<volume_id>"
}
```

```bash
$ terraform plan -generate-config-out=volume.tf
$ terraform apply
```
