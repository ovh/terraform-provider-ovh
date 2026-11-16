---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_multipart_uploads (Data Source)

List the incomplete multipart uploads of an S3&trade; compatible object storage bucket in a public cloud project. Every page of the listing is read.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

## Example Usage

```terraform
data "ovh_cloud_storage_object_bucket_multipart_uploads" "uploads" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-bucket"
}

output "multipart_upload_keys" {
  value = { for upload in data.ovh_cloud_storage_object_bucket_multipart_uploads.uploads.uploads : upload.id => upload.key }
}
```

## Argument Reference

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `bucket_id` - (Required) ID of the bucket, in the `<REGION>_<name>` form (e.g. `GRA_my-bucket`): the `id` of an `ovh_cloud_storage_object_bucket`.
* `key` - (Optional) Object key. If set, only the multipart uploads of this object are returned.

## Attributes Reference

* `uploads` - List of the incomplete multipart uploads of the bucket:
  * `id` - Unique identifier of the multipart upload.
  * `key` - Object key (path within the bucket).
  * `storage_class` - Storage class of the object (`DEEP_ARCHIVE`, `GLACIER_IR`, `HIGH_PERF`, `STANDARD`, `STANDARD_IA`).
  * `created_at` - Multipart upload initiation timestamp.
