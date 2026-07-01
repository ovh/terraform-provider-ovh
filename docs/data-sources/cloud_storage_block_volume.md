---
subcategory: "Cloud Storage"
---

# ovh_cloud_storage_block_volume (Data Source)

Get a block storage volume in a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_storage_block_volume" "volume" {
  service_name = "xxxxxxxxx"
  id           = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `id` - (Required) The ID of the volume.

## Attributes Reference

* `name` - Volume name.
* `location` - Location of the volume:
  * `region` - Region.
* `size` - Size of the volume in GB.
* `volume_type` - Volume type (`CLASSIC`, `HIGH_SPEED`, `HIGH_SPEED_GEN2`).
* `bootable` - Whether the volume is bootable.
* `status` - Volume status (`AVAILABLE`, `IN_USE`, `CREATING`, `DELETING`, `ATTACHING`, `DETACHING`, `EXTENDING`, `ERROR`, `ERROR_DELETING`, `ERROR_BACKING_UP`, `ERROR_RESTORING`, `ERROR_EXTENDING`).
* `resource_status` - Volume readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `encryption` - Encryption configuration of the volume:
  * `enabled` - Whether the volume is encrypted at rest with LUKS.
* `attached_instances` - Instances the volume is attached to:
  * `id` - Instance ID.
