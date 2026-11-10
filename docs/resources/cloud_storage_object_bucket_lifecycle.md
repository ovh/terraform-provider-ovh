---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_lifecycle

Manages the lifecycle configuration of an S3&trade; compatible object storage bucket in a public cloud project: rules that expire objects, transition them to another storage class, clean up noncurrent versions or abort incomplete multipart uploads.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

This resource owns the whole lifecycle configuration of the bucket. Creating and updating it replace every rule of the bucket with the `rules` list, and the change is applied synchronously. Destroying it removes every lifecycle rule of the bucket. If the lifecycle configuration is removed outside of Terraform, the next `terraform apply` recreates it.

## Example Usage

```terraform
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "<Public cloud project id>"
  name         = "my-versioned-bucket"
  region       = "GRA"

  versioning = {
    status = "ENABLED"
  }
}

resource "ovh_cloud_storage_object_bucket_lifecycle" "lifecycle" {
  service_name = ovh_cloud_storage_object_bucket.bucket.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.bucket.id

  rules = [
    {
      id     = "expire-logs"
      status = "ENABLED"
      filter = {
        prefix = "logs/"
      }
      expiration = {
        days = 30
      }
    },
    {
      id     = "archive-reports"
      status = "ENABLED"
      filter = {
        prefix = "reports/"
        tags = {
          retention = "long"
        }
      }
      transitions = [
        {
          days          = 90
          storage_class = "STANDARD_IA"
        },
        {
          days          = 365
          storage_class = "DEEP_ARCHIVE"
        }
      ]
    },
    {
      id     = "clean-noncurrent-versions"
      status = "ENABLED"
      noncurrent_version_transitions = [
        {
          noncurrent_days = 14
          storage_class   = "STANDARD_IA"
        }
      ]
      noncurrent_version_expiration = {
        noncurrent_days           = 60
        newer_noncurrent_versions = 3
      }
      expiration = {
        expired_object_delete_marker = true
      }
    },
    {
      id     = "abort-multipart"
      status = "ENABLED"
      abort_incomplete_multipart_upload = {
        days_after_initiation = 7
      }
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `bucket_id` - (Required) ID of the bucket, in the `<REGION>_<name>` form (e.g. `GRA_my-bucket`): the `id` of an `ovh_cloud_storage_object_bucket`. **Changing this value recreates the resource.**
* `rules` - (Required) List of lifecycle rules, at least one. Every change is applied in place:
  * `id` - (Required) Rule identifier, unique in the bucket.
  * `status` - (Required) Rule status. One of `ENABLED`, `DISABLED`.
  * `filter` - (Optional) Objects the rule applies to. When omitted, the rule applies to every object of the bucket:
    * `prefix` - (Optional) Only objects whose key starts with this prefix.
    * `tags` - (Optional) Only objects carrying every one of these tags, as a map of strings.
    * `object_size_greater_than` - (Optional) Only objects larger than this size, in bytes. Must be at least `1`.
    * `object_size_less_than` - (Optional) Only objects smaller than this size, in bytes. Must be at least `1`.
  * `expiration` - (Optional) Expiration of the current object versions:
    * `date` - (Optional) Date after which objects expire, in the `YYYY-MM-DD` form. Conflicts with `days`.
    * `days` - (Optional) Number of days after creation before objects expire. Must be at least `1`. Conflicts with `date`.
    * `expired_object_delete_marker` - (Optional) Whether delete markers with no noncurrent version left are removed.
  * `transitions` - (Optional) Transitions of the current object versions to another storage class, at least one when set:
    * `date` - (Optional) Date after which objects are transitioned, in the `YYYY-MM-DD` form. Conflicts with `days`.
    * `days` - (Optional) Number of days after creation before objects are transitioned. Must be at least `1`. Conflicts with `date`.
    * `storage_class` - (Required) Target storage class. One of `DEEP_ARCHIVE`, `GLACIER_IR`, `STANDARD`, `STANDARD_IA`.
  * `noncurrent_version_expiration` - (Optional) Expiration of the noncurrent object versions (requires versioning on the bucket):
    * `noncurrent_days` - (Required) Number of days after a version becomes noncurrent before it expires. Must be at least `1`.
    * `newer_noncurrent_versions` - (Optional) Number of most recent noncurrent versions to keep. Must be at least `1`.
  * `noncurrent_version_transitions` - (Optional) Transitions of the noncurrent object versions to another storage class (requires versioning on the bucket), at least one when set:
    * `noncurrent_days` - (Required) Number of days after a version becomes noncurrent before it is transitioned. Must be at least `1`.
    * `newer_noncurrent_versions` - (Optional) Number of most recent noncurrent versions to keep in their storage class. Must be at least `1`.
    * `storage_class` - (Required) Target storage class. One of `DEEP_ARCHIVE`, `GLACIER_IR`, `STANDARD`, `STANDARD_IA`.
  * `abort_incomplete_multipart_upload` - (Optional) Cleanup of incomplete multipart uploads:
    * `days_after_initiation` - (Required) Number of days after initiation before an incomplete multipart upload is aborted. Must be at least `1`.

## Attributes Reference

The following attributes are exported:

* `id` - Lifecycle configuration identifier, the same value as `bucket_id`.

## Import

A bucket lifecycle configuration can be imported using the `service_name` and the `bucket_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_storage_object_bucket_lifecycle.lifecycle
  id = "<service_name>/<bucket_id>"
}
```

```bash
$ terraform import ovh_cloud_storage_object_bucket_lifecycle.lifecycle service_name/bucket_id
```
