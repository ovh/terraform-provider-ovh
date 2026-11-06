---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_objects (Data Source)

List the objects stored in an S3&trade; compatible object storage bucket of a public cloud project. Every page of the listing is read.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

## Example Usage

Every object of the bucket:

```terraform
data "ovh_cloud_storage_object_bucket_objects" "objects" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-bucket"
}
```

The objects directly under the `documents/` prefix, the deeper levels grouped into one entry per sub-prefix:

```terraform
data "ovh_cloud_storage_object_bucket_objects" "documents" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-bucket"
  prefix       = "documents/"
  delimiter    = "/"
}

output "document_keys" {
  value = [
    for object in data.ovh_cloud_storage_object_bucket_objects.documents.objects :
    object.key if !object.is_common_prefix
  ]
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `bucket_id` - (Required) Bucket identifier, in the `<REGION>_<name>` form (e.g. `GRA_my-bucket`).
* `prefix` - (Optional) If set, only objects whose key starts with this prefix are returned.
* `delimiter` - (Optional) If set, the keys that share a prefix up to the next occurrence of the delimiter are grouped into one entry with `is_common_prefix = true`. When unset, every object is listed recursively.
* `with_versions` - (Optional) Whether to list every version and delete marker of the objects, not only their latest version.

## Attributes Reference

* `objects` - List of objects:
  * `key` - Object key, its full path within the bucket (e.g. `documents/report.pdf`). For a common prefix, the prefix itself (e.g. `documents/2026/`).
  * `size` - Object size, in bytes.
  * `etag` - Entity tag of the object.
  * `last_modified_at` - Last modification date of the object.
  * `storage_class` - Storage class of the object (`DEEP_ARCHIVE`, `GLACIER_IR`, `HIGH_PERF`, `STANDARD`, `STANDARD_IA`).
  * `version_id` - Version ID of the object.
  * `is_latest` - Whether this version is the latest version of the object.
  * `is_delete_marker` - Whether this version is a delete marker.
  * `is_common_prefix` - Whether this entry groups the keys sharing a prefix (only with `delimiter`) rather than describing an object.
  * `legal_hold` - Legal hold status of the object (`ON`, `OFF`).
  * `lock` - Object lock retention of the object:
    * `mode` - Retention mode (`COMPLIANCE`, `GOVERNANCE`).
    * `retain_until` - Date until which the object is retained.
  * `restore_status` - Restore status of an archived object:
    * `in_progress` - Whether a restore is in progress.
    * `expires_at` - Date when the restored copy expires.
  * `replication_status` - Replication status of the object (`COMPLETED`, `FAILED`, `NONE`, `PENDING`, `REPLICA`).

Attributes the API does not return for an entry are `null`.
