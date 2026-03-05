---
subcategory : "Cloud Storage"
---

# ovh_cloud_storage_file_share

Creates a file storage share (NFS) in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_storage_file_share" "share" {
  service_name = "xxxxxxxxxx"
  name         = "my-share"
  size         = 100
  region       = "GRA1"
  protocol     = "NFS"
  share_type   = "STANDARD_1AZ"
  network_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subnet_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  description  = "My NFS share"

  access_rules {
    access_to    = "10.0.0.0/24"
    access_level = "READ_WRITE"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `name` - (Required) File share name.
* `size` - (Required) Size of the file share in GB.
* `region` - (Required) Region where the file share will be created. **Changing this value recreates the resource.**
* `protocol` - (Required) File share protocol (`NFS`). **Changing this value recreates the resource.**
* `share_type` - (Required) File share type (e.g. `STANDARD_1AZ`). **Changing this value recreates the resource.**
* `network_id` - (Required) Network ID to attach the file share to. **Changing this value recreates the resource.**
* `subnet_id` - (Required) Subnet ID to attach the file share to. **Changing this value recreates the resource.**
* `availability_zone` - (Optional) Availability zone where the file share will be created. **Changing this value recreates the resource.**
* `description` - (Optional) File share description.
* `access_rules` - (Optional) Access rules for the file share. Each rule has:
  * `access_to` - (Required) IP address or CIDR to grant access to.
  * `access_level` - (Required) Access level (`READ_WRITE`, `READ_ONLY`).

## Attributes Reference

The following attributes are exported:

* `id` - File share ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the file share.
* `updated_at` - Last update date of the file share.
* `resource_status` - File share readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the file storage share:
  * `location` - Current location:
    * `region` - Region.
    * `availability_zone` - Availability zone.
  * `name` - File share name.
  * `description` - File share description.
  * `size` - Size of the file share in GB.
  * `protocol` - File share protocol.
  * `share_type` - File share type.
  * `network_id` - Network ID.
  * `subnet_id` - Subnet ID.
  * `status` - File share status.
  * `export_locations` - Export locations for the file share:
    * `path` - Export path.
    * `preferred` - Whether this is the preferred export location.
  * `access_rules` - Current access rules for the file share:
    * `id` - Access rule ID.
    * `access_to` - IP address or CIDR.
    * `access_level` - Access level.
    * `state` - Access rule state.
    * `created_at` - Access rule creation date.

## Import

A cloud storage file share can be imported using the `service_name` and `file_share_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_storage_file_share.share
  id = "<service_name>/<file_share_id>"
}
```

```bash
$ terraform import ovh_cloud_storage_file_share.share service_name/file_share_id
```
