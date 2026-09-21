---
subcategory : "Object Storage"
---

# ovh_cloud_s3_bucket

Creates an S3&trade; compatible object storage bucket in a public cloud project.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

## Example Usage

```terraform
resource "ovh_cloud_s3_bucket" "bucket" {
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

Object lock (WORM) requires versioning to be enabled:

```terraform
resource "ovh_cloud_s3_bucket" "locked" {
  service_name = <Public cloud project id>
  name         = "my-worm-bucket"
  region       = "GRA"

  versioning = {
    status = "ENABLED"
  }

  object_lock = {
    mode           = "GOVERNANCE"
    retention_days = 30
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `name` - (Required) Bucket name (must be globally unique and DNS-compatible). **Changing this value recreates the resource.**
* `region` - (Required) Region identifier where the bucket will be created (e.g. `GRA`, `SBG`, `BHS`). **Changing this value recreates the resource.**
* `owner_user_id` - (Optional) Owner user identifier.
* `tags` - (Optional) Metadata tags for the bucket, as a map of strings.
* `encryption` - (Optional) Server-side encryption configuration:
  * `algorithm` - (Required) Encryption algorithm. One of `AES256`, `PLAINTEXT`.
* `versioning` - (Optional) Versioning configuration:
  * `status` - (Required) Versioning status. One of `DISABLED`, `ENABLED`, `SUSPENDED`.
* `object_lock` - (Optional) Object lock (WORM) configuration. Requires `versioning.status` to be `ENABLED`:
  * `mode` - (Required) Object lock retention mode. One of `COMPLIANCE`, `GOVERNANCE`.
  * `retention_days` - (Required) Number of days to retain objects.

## Attributes Reference

The following attributes are exported:

* `id` - Bucket identifier.
* `object_lock`:
  * `retention_years` - Number of years to retain objects. Read-only alternative to `retention_days`, set only on buckets locked outside of Terraform.
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
  * `object_lock` - Current object lock configuration:
    * `mode` - Object lock retention mode.
    * `retention_days` - Number of days to retain objects.
    * `retention_years` - Number of years to retain objects.
  * `tags` - Current metadata tags.
  * `virtual_host` - Bucket virtual host, as a hostname without scheme (for example `my-data-bucket.s3.gra.io.cloud.ovh.net`). Only returned on a single bucket read.
  * `objects_count` - Bucket total objects count. Only returned on a single bucket read.
  * `objects_size` - Bucket total objects size in bytes. Only returned on a single bucket read.

## Import

An S3&trade; compatible bucket can be imported using the `service_name` and the bucket `id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_s3_bucket.bucket
  id = "<service_name>/<bucket_id>"
}
```

```bash
$ terraform import ovh_cloud_s3_bucket.bucket service_name/bucket_id
```
