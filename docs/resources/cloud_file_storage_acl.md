---
subcategory : "Cloud Storage"
---

# ovh_cloud_file_storage_acl

Creates an access rule (ACL) on a public cloud file storage share, controlling which IP addresses can access it. The region is inherited from the parent share.

## Example Usage

```terraform
resource "ovh_cloud_storage_file_share" "share" {
  service_name = "<Public cloud project id>"
  name         = "my-share"
  size         = 150
  region       = "GRA1"
  protocol     = "NFS"
  share_type   = "STANDARD_1AZ"
}

resource "ovh_cloud_file_storage_acl" "acl" {
  service_name = ovh_cloud_storage_file_share.share.service_name
  share_id     = ovh_cloud_storage_file_share.share.id
  access_to    = "10.0.0.0/24"
  access_level = "READ_WRITE"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `share_id` - (Required) ID of the file storage share the access rule applies to. **Changing this value recreates the resource.**
* `access_to` - (Required) IP address or CIDR allowed to access the file storage share. **Changing this value recreates the resource.**
* `access_level` - (Required) Access level granted (`READ_WRITE`, `READ_ONLY`).

## Attributes Reference

The following attributes are exported:

* `id` - Access rule ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the access rule.
* `updated_at` - Last update date of the access rule.
* `resource_status` - Access rule readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current observed state of the access rule:
  * `access_to` - IP address or CIDR allowed to access the file storage share.
  * `access_level` - Access level granted.
  * `state` - Current state of the access rule (`ACTIVE`, `APPLYING`, `DENYING`, `ERROR`).
  * `created_at` - Creation date of the access rule.

## Import

A cloud file storage access rule can be imported using the `service_name`, `share_id` and `acl_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_file_storage_acl.acl
  id = "<service_name>/<share_id>/<acl_id>"
}
```

```bash
$ terraform import ovh_cloud_file_storage_acl.acl service_name/share_id/acl_id
```
