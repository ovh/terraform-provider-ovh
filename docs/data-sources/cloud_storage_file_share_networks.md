---
subcategory: "Cloud Storage"
---

# ovh_cloud_storage_file_share_networks (Data Source)

List the file storage share networks in a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_storage_file_share_networks" "networks" {
  service_name = "<public cloud project ID>"
}
```

Filter the share networks by region:

```hcl
data "ovh_cloud_storage_file_share_networks" "networks" {
  service_name = "<public cloud project ID>"
  region       = "GRA9"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `region` - (Optional) If set, only share networks located in this region are returned.

## Attributes Reference

* `share_networks` - List of share networks:
  * `id` - Share network ID.
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
