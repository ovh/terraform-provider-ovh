---
subcategory: "Block Storage"
---

# ovh_cloud_storage_block_volume_snapshot

Creates a block storage volume snapshot in a public cloud project.

## Example Usage

```hcl
resource "ovh_cloud_storage_block_volume_snapshot" "snapshot" {
  service_name = "xxxxxxxxx"
  name         = "my-snapshot"
  description  = "Snapshot before upgrade"
  region       = "GRA9"
  volume_id    = ovh_cloud_storage_block_volume.volume.id
}
```

### Restore a volume from a snapshot

```hcl
resource "ovh_cloud_storage_block_volume" "restored" {
  service_name = "xxxxxxxxx"
  name         = "restored-volume"
  size         = 10
  region       = "GRA9"
  volume_type  = "CLASSIC"

  create_from = {
    snapshot_id = ovh_cloud_storage_block_volume_snapshot.snapshot.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the public cloud project. Changing this value recreates the resource.
* `name` - (Required) The name of the snapshot.
* `description` - (Optional) A description for the snapshot.
* `region` - (Required) The region where the snapshot will be created. Changing this value recreates the resource.
* `volume_id` - (Required) The ID of the volume to snapshot. Changing this value recreates the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The snapshot ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the snapshot.
* `updated_at` - Last update date of the snapshot.
* `resource_status` - Snapshot readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the snapshot:
  * `location` - Current location:
    * `region` - Region.
  * `name` - Snapshot name.
  * `description` - Snapshot description.
  * `volume_id` - ID of the snapshotted volume.
  * `size` - Size of the snapshot in GB.

## Import

A block storage volume snapshot can be imported using the `service_name` and `id` separated by a `/`:

```bash
terraform import ovh_cloud_storage_block_volume_snapshot.snapshot service_name/snapshot_id
```
