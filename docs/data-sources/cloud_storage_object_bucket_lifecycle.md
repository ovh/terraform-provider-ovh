---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_lifecycle (Data Source)

Get the lifecycle configuration of an S3&trade; compatible object storage bucket in a public cloud project.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

A bucket without lifecycle configuration is not an error: `rules` is an empty list.

## Example Usage

```terraform
data "ovh_cloud_storage_object_bucket_lifecycle" "lifecycle" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-bucket"
}

output "lifecycle_rule_ids" {
  value = [for rule in data.ovh_cloud_storage_object_bucket_lifecycle.lifecycle.rules : rule.id]
}
```

## Argument Reference

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `bucket_id` - (Required) ID of the bucket, in the `<REGION>_<name>` form (e.g. `GRA_my-bucket`): the `id` of an `ovh_cloud_storage_object_bucket`.

## Attributes Reference

* `id` - Lifecycle configuration identifier, the same value as `bucket_id`.
* `rules` - List of lifecycle rules, empty when the bucket has no lifecycle configuration:
  * `id` - Rule identifier.
  * `status` - Rule status (`ENABLED`, `DISABLED`).
  * `filter` - Objects the rule applies to, `null` when the rule applies to every object of the bucket:
    * `prefix` - Only objects whose key starts with this prefix.
    * `tags` - Only objects carrying every one of these tags, as a map of strings.
    * `object_size_greater_than` - Only objects larger than this size, in bytes.
    * `object_size_less_than` - Only objects smaller than this size, in bytes.
  * `expiration` - Expiration of the current object versions:
    * `date` - Date after which objects expire, in the `YYYY-MM-DD` form.
    * `days` - Number of days after creation before objects expire.
    * `expired_object_delete_marker` - Whether delete markers with no noncurrent version left are removed.
  * `transitions` - Transitions of the current object versions to another storage class:
    * `date` - Date after which objects are transitioned, in the `YYYY-MM-DD` form.
    * `days` - Number of days after creation before objects are transitioned.
    * `storage_class` - Target storage class (`DEEP_ARCHIVE`, `GLACIER_IR`, `STANDARD`, `STANDARD_IA`).
  * `noncurrent_version_expiration` - Expiration of the noncurrent object versions:
    * `noncurrent_days` - Number of days after a version becomes noncurrent before it expires.
    * `newer_noncurrent_versions` - Number of most recent noncurrent versions to keep.
  * `noncurrent_version_transitions` - Transitions of the noncurrent object versions to another storage class:
    * `noncurrent_days` - Number of days after a version becomes noncurrent before it is transitioned.
    * `newer_noncurrent_versions` - Number of most recent noncurrent versions to keep in their storage class.
    * `storage_class` - Target storage class (`DEEP_ARCHIVE`, `GLACIER_IR`, `HIGH_PERF`, `STANDARD`, `STANDARD_IA`).
  * `abort_incomplete_multipart_upload` - Cleanup of incomplete multipart uploads:
    * `days_after_initiation` - Number of days after initiation before an incomplete multipart upload is aborted.

Optional attributes are `null` when the API does not return them.
