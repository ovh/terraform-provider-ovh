---
subcategory: "Cloud Storage"
---

# ovh_cloud_storage_file_share_snapshot (Data Source)

Get a file storage snapshot (NFS share snapshot) in a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_storage_file_share_snapshot" "snapshot" {
  service_name = "<public cloud project ID>"
  id           = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `id` - (Required) The ID of the snapshot.

## Attributes Reference

* `name` - Snapshot name.
* `description` - Snapshot description.
* `share_id` - ID of the snapshotted file share.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the snapshot.
* `updated_at` - Last update date of the snapshot.
* `resource_status` - Snapshot readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the file storage snapshot:
  * `name` - Snapshot name.
  * `description` - Snapshot description.
  * `share_id` - ID of the snapshotted file share.
  * `size` - Size of the snapshot in GB.
  * `location` - Current location:
    * `region` - Region.
    * `availability_zone` - Availability zone.
