---
subcategory : "Cloud Storage"
---

# ovh_cloud_storage_file_share_snapshot

Creates a snapshot of a file storage share (NFS) in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_storage_file_share" "share" {
  service_name = "xxxxxxxxxx"
  name         = "my-share"
  size         = 100
  region       = "GRA1"
  protocol     = "NFS"
  share_type   = "STANDARD_1AZ"
  network_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subnet_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

resource "ovh_cloud_storage_file_share_snapshot" "snapshot" {
  service_name = ovh_cloud_storage_file_share.share.service_name
  region       = ovh_cloud_storage_file_share.share.region
  share_id     = ovh_cloud_storage_file_share.share.id
  name         = "my-snapshot"
  description  = "Daily snapshot of my-share"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `region` - (Required) Region where the snapshot will be created. **Changing this value recreates the resource.**
* `share_id` - (Required) ID of the file share to snapshot. **Changing this value recreates the resource.**
* `name` - (Optional) Snapshot name.
* `description` - (Optional) Snapshot description.

## Attributes Reference

The following attributes are exported:

* `id` - Snapshot ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the snapshot.
* `updated_at` - Last update date of the snapshot.
* `resource_status` - Snapshot readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the file storage snapshot:
  * `location` - Current location:
    * `region` - Region.
  * `name` - Snapshot name.
  * `description` - Snapshot description.
  * `share_id` - ID of the snapshotted file share.
  * `snapshot_size` - Size of the snapshot in GB.
  * `share_size` - Size of the source file share in GB.
  * `share_proto` - Protocol of the source file share (e.g. NFS).

## Import

A cloud storage file share snapshot can be imported using the `service_name` and `snapshot_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_storage_file_share_snapshot.snapshot
  id = "<service_name>/<snapshot_id>"
}
```

```bash
$ terraform import ovh_cloud_storage_file_share_snapshot.snapshot service_name/snapshot_id
```
