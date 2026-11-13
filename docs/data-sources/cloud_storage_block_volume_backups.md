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
  * `status` - Backup status (`AVAILABLE`, `CREATING`, `DELETED`, `DELETING`, `ERROR`, `ERROR_DELETING`, `ERROR_RESTORING`, `RESTORING`, `UNMANAGING`).
  * `volume_type` - Type of the backed-up volume (`CLASSIC`, `HIGH_SPEED`, `HIGH_SPEED_GEN2`). Null when the backup does not record it.
  * `encryption` - Encryption configuration of the backed-up volume. Null when the backup does not record it:
    * `enabled` - Whether the backed-up volume is encrypted at rest with LUKS.
    * `kms` - Customer-managed key (CMK) reference used to encrypt the backed-up volume:
      * `domain_id` - OKMS domain ID owning the service key.
      * `service_key_id` - OKMS service key ID used to encrypt the backed-up volume.
  * `resource_status` - Backup readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
