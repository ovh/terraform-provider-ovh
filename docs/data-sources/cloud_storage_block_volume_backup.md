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
* `status` - Backup status (`AVAILABLE`, `CREATING`, `DELETED`, `DELETING`, `ERROR`, `ERROR_DELETING`, `ERROR_RESTORING`, `RESTORING`, `UNMANAGING`).
* `volume_type` - Type of the backed-up volume (`CLASSIC`, `HIGH_SPEED`, `HIGH_SPEED_GEN2`). Null when the backup does not record it.
* `encryption` - Encryption configuration of the backed-up volume. Null when the backup does not record it:
  * `enabled` - Whether the backed-up volume is encrypted at rest with LUKS.
  * `kms` - Customer-managed key (CMK) reference used to encrypt the backed-up volume:
    * `domain_id` - OKMS domain ID owning the service key.
    * `service_key_id` - OKMS service key ID used to encrypt the backed-up volume.
* `resource_status` - Backup readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
