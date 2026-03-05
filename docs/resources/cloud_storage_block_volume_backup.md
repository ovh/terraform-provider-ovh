---
subcategory: "Cloud Storage"
---

# ovh_cloud_storage_block_volume_backup

Creates a block storage volume backup in a public cloud project.

## Example Usage

```hcl
resource "ovh_cloud_storage_block_volume_backup" "backup" {
  service_name = "xxxxxxxxx"
  name         = "my-backup"
  description  = "Daily backup"
  region       = "GRA9"
  volume_id    = ovh_cloud_storage_block_volume.volume.id
}
```

### Restore a volume from a backup

```hcl
resource "ovh_cloud_storage_block_volume" "restored" {
  service_name = "xxxxxxxxx"
  name         = "restored-volume"
  size         = 10
  region       = "GRA9"
  volume_type  = "CLASSIC"

  create_from {
    backup_id = ovh_cloud_storage_block_volume_backup.backup.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the public cloud project. Changing this value recreates the resource.
* `name` - (Required) The name of the backup.
* `description` - (Optional) A description for the backup.
* `region` - (Required) The region where the backup will be created. Changing this value recreates the resource.
* `volume_id` - (Required) The ID of the volume to back up. Changing this value recreates the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The backup ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the backup.
* `updated_at` - Last update date of the backup.
* `resource_status` - Backup readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the backup:
  * `location` - Current location:
    * `region` - Region.
  * `name` - Backup name.
  * `description` - Backup description.
  * `volume_id` - ID of the backed-up volume.
  * `size` - Size of the backup in GB.

## Import

A block storage volume backup can be imported using the `service_name` and `id` separated by a `/`:

```bash
terraform import ovh_cloud_storage_block_volume_backup.backup service_name/backup_id
```
