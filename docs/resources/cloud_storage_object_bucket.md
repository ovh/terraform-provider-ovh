---
subcategory: "Object Storage"
---

# ovh_cloud_storage_object_bucket

Creates an S3&trade; compatible object storage bucket in a public cloud project.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

## Example Usage

```terraform
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name  = <Public cloud project id>
  name          = "my-data-bucket"
  region        = "GRA"
  owner_user_id = "<owner user id>"

  versioning = {
    status = "ENABLED"
  }

  encryption = {
    algorithm = "AES256"
  }

  tags = {
    env = "production"
  }
}
```

Object lock (WORM) requires `versioning.status` to be `ENABLED`, and can only be armed when the bucket is created. The default retention rule (`mode` + `retention_days`) is optional:

```terraform
resource "ovh_cloud_storage_object_bucket" "locked" {
  service_name = <Public cloud project id>
  name         = "my-worm-bucket"
  region       = "GRA"

  versioning = {
    status = "ENABLED"
  }

  object_lock = {
    status         = "ENABLED"
    mode           = "GOVERNANCE"
    retention_days = 30
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `name` - (Required) Bucket name (must be globally unique and DNS-compatible): 3 to 63 characters, lowercase letters, digits, dots and hyphens only, starting and ending with a letter or a digit. **Changing this value recreates the resource.**
* `region` - (Required) Region identifier where the bucket will be created. Must be upper-case (e.g. `GRA`, `SBG`, `BHS`). **Changing this value recreates the resource.**
* `owner_user_id` - (Optional) Owner user identifier.
* `tags` - (Optional) Metadata tags for the bucket, as a map of strings.
* `encryption` - (Optional) Server-side encryption configuration:
  * `algorithm` - (Required) Encryption algorithm. One of `AES256`, `PLAINTEXT`.
* `versioning` - (Optional) Versioning configuration:
  * `status` - (Required) Versioning status. One of `DISABLED`, `ENABLED`, `SUSPENDED`.
* `object_lock` - (Optional) Object lock (WORM) configuration. Object lock can only be armed at bucket creation. **Adding, changing or removing this block recreates the resource.**
  * `status` - (Required) Whether object lock is armed on the bucket. One of `ENABLED`, `DISABLED`. `ENABLED` requires `versioning` to be set with `status = "ENABLED"`.
  * `mode` - (Optional) Default retention mode. One of `COMPLIANCE`, `GOVERNANCE`. Must be set together with `retention_days`, and only when `status` is `ENABLED`.
  * `retention_days` - (Optional) Number of days the default retention rule locks new objects for. Must be at least `1`, and set together with `mode`.

## Attributes Reference

The following attributes are exported:

* `id` - Bucket identifier.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the bucket.
* `updated_at` - Last update date of the bucket.
* `resource_status` - Bucket readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `SUSPENDED`, `UNKNOWN`, `UPDATING`).
* `current_state` - Current observed state of the bucket:
  * `name` - Bucket name.
  * `location` - Geographic region where the bucket is located:
    * `region` - Region identifier.
  * `encryption` - Current encryption configuration:
    * `algorithm` - Encryption algorithm.
  * `versioning` - Current versioning configuration:
    * `status` - Versioning status.
  * `object_lock` - Current object lock configuration, always present:
    * `status` - Whether object lock is armed on the bucket (`ENABLED`, `DISABLED`).
    * `mode` - Default retention mode, `null` without a default retention rule.
    * `retention_days` - Number of days the default retention rule locks new objects for, `null` without a default retention rule.
    * `retention_years` - Number of years to retain objects. Read-only alternative to `retention_days`, set only on buckets locked outside of Terraform.
  * `tags` - Current metadata tags.
  * `virtual_host` - Bucket virtual host, as a hostname without scheme (for example `my-data-bucket.s3.gra.io.cloud.ovh.net`). Only returned on a single bucket read.
  * `objects_count` - Bucket total objects count. Only returned on a single bucket read.
  * `objects_size` - Bucket total objects size in bytes. Only returned on a single bucket read.

## Import

An S3&trade; compatible bucket can be imported using the `service_name` and the bucket `id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_storage_object_bucket.bucket
  id = "<service_name>/<bucket_id>"
}
```

~> __NOTE__ The bucket `id` is the bare bucket name on a single-region API instance, and `<REGION>_<name>` (for example `GRA_my-data-bucket`) on a multi-region one.

```bash
$ terraform import ovh_cloud_storage_object_bucket.bucket service_name/bucket_id
```
