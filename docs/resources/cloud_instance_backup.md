---
subcategory: "Instances"
---

# ovh_cloud_instance_backup

Creates an on-demand backup of an instance in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_instance_backup" "backup" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  instance_id  = "<instance id>"
  name         = "my-instance-backup"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. Changing this value recreates the resource.
* `region` - (Required) Region where the backup will be created. Changing this value recreates the resource.
* `instance_id` - (Required) ID of the instance to back up. Changing this value recreates the resource.
* `name` - (Required) Backup name. Changing this value recreates the resource.

-> All attributes are immutable: the resource does not support in-place updates, any change requires replacement.

## Attributes Reference

The following attributes are exported:

* `id` - Backup ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the backup.
* `updated_at` - Last update date of the backup.
* `resource_status` - Backup readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`).
* `current_state` - Current state of the instance backup:
  * `instance` - Source instance reference:
    * `id` - Instance unique identifier.
  * `location` - Current location:
    * `region` - Region.
  * `min_disk` - Minimum disk size in GB required to boot.
  * `min_ram` - Minimum RAM in MB required to boot.
  * `name` - Current backup name.
  * `size` - Image size in bytes.
  * `status` - Image status in the backend.
  * `visibility` - Image visibility.

## Import

An instance backup can be imported using the `service_name` and the `id` of the backup, separated by "/":

```bash
terraform import ovh_cloud_instance_backup.backup service_name/backup_id
```
