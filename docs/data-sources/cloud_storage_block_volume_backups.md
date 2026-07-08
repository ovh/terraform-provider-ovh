---
subcategory: "Block Storage"
---

# ovh_cloud_storage_block_volume_backups (Data Source)

List the backups of a block storage volume in a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_storage_block_volume_backups" "backups" {
  service_name = "xxxxxxxxx"
  region       = "GRA9"
  volume_id    = ovh_cloud_storage_block_volume.volume.id
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `region` - (Required) The region where the backups reside.
* `volume_id` - (Required) The ID of the volume whose backups to list.

## Attributes Reference

* `backups` - List of backups for the volume:
  * `id` - Backup ID.
  * `name` - Backup name.
  * `description` - Backup description.
  * `location` - Location of the backup:
    * `region` - Region.
  * `volume_id` - ID of the backed-up volume.
  * `size` - Size of the backup in GB.
  * `resource_status` - Backup readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
