---
subcategory: "Block Storage"
---

# ovh_cloud_storage_block_volume_snapshots (Data Source)

List the snapshots of a block storage volume in a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_storage_block_volume_snapshots" "snapshots" {
  service_name = "xxxxxxxxx"
  region       = "GRA9"
  volume_id    = ovh_cloud_storage_block_volume.volume.id
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `region` - (Required) The region where the snapshots reside.
* `volume_id` - (Required) The ID of the volume whose snapshots to list.

## Attributes Reference

* `snapshots` - List of snapshots for the volume:
  * `id` - Snapshot ID.
  * `name` - Snapshot name.
  * `description` - Snapshot description.
  * `location` - Location of the snapshot:
    * `region` - Region.
  * `volume_id` - ID of the snapshotted volume.
  * `size` - Size of the snapshot in GB.
  * `resource_status` - Snapshot readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
