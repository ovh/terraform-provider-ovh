---
subcategory: "Cloud Storage"
---

# ovh_cloud_storage_file_share_network (Data Source)

Get a file storage share network in a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_storage_file_share_network" "network" {
  service_name = "<public cloud project ID>"
  id           = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `id` - (Required) The ID of the share network.

## Attributes Reference

* `name` - Share network name.
* `description` - Share network description.
* `network_id` - ID of the private network.
* `subnet_id` - ID of the subnet.
* `location` - Location of the share network:
  * `region` - Region where the share network resides.
  * `availability_zone` - Availability zone where the share network resides.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the share network.
* `updated_at` - Last update date of the share network.
* `resource_status` - Share network readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the file storage share network:
  * `name` - Share network name.
  * `description` - Share network description.
  * `network_id` - ID of the private network.
  * `subnet_id` - ID of the subnet.
  * `location` - Current location:
    * `region` - Region.
    * `availability_zone` - Availability zone.
