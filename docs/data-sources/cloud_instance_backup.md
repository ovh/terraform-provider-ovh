---
subcategory: "Instances"
---

# ovh_cloud_instance_backup (Data Source)

Use this data source to retrieve information about an instance backup in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_instance_backup" "backup" {
  service_name = "<Public cloud project id>"
  id           = "<backup id>"
}

output "backup_status" {
  value = data.ovh_cloud_instance_backup.backup.resource_status
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `id` - (Required) Backup ID.

## Attributes Reference

The following attributes are exported:

* `name` - Backup name.
* `location` - Location of the backup:
  * `region` - Region.
* `instance_id` - ID of the backed-up instance.
* `min_disk` - Minimum disk size in GB required to boot.
* `min_ram` - Minimum RAM in MB required to boot.
* `size` - Image size in bytes.
* `status` - Image status in the backend.
* `visibility` - Image visibility.
* `resource_status` - Backup readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`).
