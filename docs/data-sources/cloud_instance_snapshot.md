---
subcategory: "Instances"
---

# ovh_cloud_instance_snapshot (Data Source)

Use this data source to retrieve information about an instance snapshot in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_instance_snapshot" "snapshot" {
  service_name = "<Public cloud project id>"
  id           = "<snapshot id>"
}

output "snapshot_status" {
  value = data.ovh_cloud_instance_snapshot.snapshot.resource_status
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `id` - (Required) Snapshot ID.

## Attributes Reference

The following attributes are exported:

* `name` - Snapshot name.
* `location` - Location of the snapshot:
  * `region` - Region.
* `instance_id` - ID of the snapshotted instance.
* `min_disk` - Minimum disk size in GB required to boot.
* `min_ram` - Minimum RAM in MB required to boot.
* `size` - Image size in bytes.
* `status` - Image status in the backend.
* `visibility` - Image visibility.
* `resource_status` - Snapshot readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`).
