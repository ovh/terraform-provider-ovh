---
subcategory: "File Storage"
---

# ovh_cloud_storage_file_share_acl (Data Source)

Get an access rule (ACL) of a public cloud file storage share.

## Example Usage

```hcl
data "ovh_cloud_storage_file_share_acl" "acl" {
  service_name = "<public cloud project ID>"
  share_id     = "00000000-0000-0000-0000-000000000000"
  id           = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `share_id` - (Required) The ID of the file storage share the access rule applies to.
* `id` - (Required) The ID of the access rule.

## Attributes Reference

* `access_to` - IP address or CIDR allowed to access the file storage share.
* `access_level` - Access level granted (`READ_WRITE`, `READ_ONLY`).
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the access rule.
* `updated_at` - Last update date of the access rule.
* `resource_status` - Access rule readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current observed state of the access rule from the infrastructure:
  * `access_to` - IP address or CIDR allowed to access the file storage share.
  * `access_level` - Access level granted.
  * `state` - Current state of the access rule (`ACTIVE`, `APPLYING`, `DENYING`, `ERROR`).
  * `created_at` - Creation date of the access rule.
