---
subcategory: "Cloud Storage"
---

# ovh_cloud_storage_block_volume_snapshot (Data Source)

Get a snapshot of a block storage volume in a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_storage_block_volume_snapshot" "snapshot" {
  service_name = "xxxxxxxxx"
  id           = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `id` - (Required) The ID of the snapshot.

## Attributes Reference

* `name` - Snapshot name.
* `description` - Snapshot description.
* `location` - Location of the snapshot:
  * `region` - Region.
* `volume_id` - ID of the snapshotted volume.
* `size` - Size of the snapshot in GB.
* `resource_status` - Snapshot readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
