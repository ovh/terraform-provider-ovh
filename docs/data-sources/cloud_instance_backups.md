---
subcategory: "Instances"
---

# ovh_cloud_instance_backups (Data Source)

Use this data source to list the backups of an instance in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_instance_backups" "backups" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  instance_id  = "<instance id>"
}

output "backup_ids" {
  value = [for backup in data.ovh_cloud_instance_backups.backups.backups : backup.id]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `region` - (Required) Region where the instance backups reside.
* `instance_id` - (Required) ID of the instance whose backups to list.

## Attributes Reference

The following attributes are exported:

* `backups` - List of backups for the instance:
  * `id` - Backup ID.
  * `name` - Backup name.
  * `location` - Location of the backup:
    * `region` - Region.
  * `instance_id` - ID of the backed-up instance.
  * `size` - Image size in bytes.
  * `status` - Image status in the backend.
  * `visibility` - Image visibility.
  * `resource_status` - Backup readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`).
