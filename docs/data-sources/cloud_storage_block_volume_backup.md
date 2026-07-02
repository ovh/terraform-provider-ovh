---
subcategory: "Block Storage"
---

# ovh_cloud_storage_block_volume_backup (Data Source)

Get a backup of a block storage volume in a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_storage_block_volume_backup" "backup" {
  service_name = "xxxxxxxxx"
  id           = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `id` - (Required) The ID of the backup.

## Attributes Reference

* `name` - Backup name.
* `description` - Backup description.
* `location` - Location of the backup:
  * `region` - Region.
* `volume_id` - ID of the backed-up volume.
* `size` - Size of the backup in GB.
* `resource_status` - Backup readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
