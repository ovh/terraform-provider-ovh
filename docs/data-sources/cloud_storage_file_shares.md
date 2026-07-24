---
subcategory: "Cloud Storage"
---

# ovh_cloud_storage_file_shares (Data Source)

List the file storage shares (NFS) in a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_storage_file_shares" "shares" {
  service_name = "<public cloud project ID>"
}
```

Filter the shares by region:

```hcl
data "ovh_cloud_storage_file_shares" "shares" {
  service_name = "<public cloud project ID>"
  region       = "GRA9"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `region` - (Optional) If set, only file shares located in this region are returned.

## Attributes Reference

* `file_shares` - List of file shares:
  * `id` - File share ID.
  * `name` - File share name.
  * `description` - File share description.
  * `size` - Size of the file share in GB.
  * `protocol` - File share protocol (`NFS`).
  * `share_type` - File share type (e.g. `STANDARD_1AZ`).
  * `location` - Location of the file share:
    * `region` - Region where the file share resides.
    * `availability_zone` - Availability zone where the file share resides.
  * `share_network_id` - ID of the share network the file share is attached to.
  * `checksum` - Computed hash representing the current target specification value.
  * `created_at` - Creation date of the file share.
  * `updated_at` - Last update date of the file share.
  * `resource_status` - File share readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
  * `current_state` - Current state of the file storage share:
    * `name` - File share name.
    * `description` - File share description.
    * `size` - Size of the file share in GB.
    * `protocol` - File share protocol.
    * `share_type` - File share type.
    * `share_network_id` - ID of the share network the file share is attached to.
    * `location` - Current location:
      * `region` - Region.
      * `availability_zone` - Availability zone.
    * `export_locations` - Export locations for the file share:
      * `path` - Export path.
      * `preferred` - Whether this is the preferred export location.
    * `capabilities` - Action-availability flags derived from the file share status:
      * `name` - Capability name.
      * `enabled` - Whether the capability is enabled.
      * `reason` - Reason why the capability is disabled.
