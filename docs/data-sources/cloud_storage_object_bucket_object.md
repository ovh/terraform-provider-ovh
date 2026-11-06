---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_object (Data Source)

Get the metadata of an object stored in an S3&trade; compatible object storage bucket of a public cloud project. The object content is not read.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

## Example Usage

```terraform
data "ovh_cloud_storage_object_bucket_object" "report" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-bucket"
  key          = "documents/report.pdf"
}

output "report_size" {
  value = data.ovh_cloud_storage_object_bucket_object.report.size
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `bucket_id` - (Required) Bucket identifier, in the `<REGION>_<name>` form (e.g. `GRA_my-bucket`).
* `key` - (Required) Object key, its full path within the bucket (e.g. `documents/report.pdf`).

## Attributes Reference

* `size` - Object size, in bytes.
* `etag` - Entity tag of the object.
* `last_modified_at` - Last modification date of the object.
* `storage_class` - Storage class of the object (`DEEP_ARCHIVE`, `GLACIER_IR`, `HIGH_PERF`, `STANDARD`, `STANDARD_IA`).
* `version_id` - Version ID of the object.
* `is_latest` - Whether this version is the latest version of the object.
* `is_delete_marker` - Whether this version is a delete marker.
* `legal_hold` - Legal hold status of the object (`ON`, `OFF`).
* `lock` - Object lock retention of the object:
  * `mode` - Retention mode (`COMPLIANCE`, `GOVERNANCE`).
  * `retain_until` - Date until which the object is retained.
* `restore_status` - Restore status of an archived object:
  * `in_progress` - Whether a restore is in progress.
  * `expires_at` - Date when the restored copy expires.
* `replication_status` - Replication status of the object (`COMPLETED`, `FAILED`, `NONE`, `PENDING`, `REPLICA`).

`legal_hold`, `lock`, `restore_status`, `replication_status` and the other optional attributes are `null` when the API does not return them.
