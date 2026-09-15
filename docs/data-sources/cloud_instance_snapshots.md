---
subcategory: "Instances"
---

# ovh_cloud_instance_snapshots (Data Source)

Use this data source to list the snapshots of an instance in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_instance_snapshots" "snapshots" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  instance_id  = "<instance id>"
}

output "snapshot_ids" {
  value = [for snapshot in data.ovh_cloud_instance_snapshots.snapshots.snapshots : snapshot.id]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `region` - (Required) Region where the instance snapshots reside.
* `instance_id` - (Required) ID of the instance whose snapshots to list.

## Attributes Reference

The following attributes are exported:

* `snapshots` - List of snapshots for the instance:
  * `id` - Snapshot ID.
  * `name` - Snapshot name.
  * `location` - Location of the snapshot:
    * `region` - Region.
  * `instance_id` - ID of the snapshotted instance.
  * `size` - Image size in bytes.
  * `status` - Image status in the backend.
  * `visibility` - Image visibility.
  * `resource_status` - Snapshot readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`).
