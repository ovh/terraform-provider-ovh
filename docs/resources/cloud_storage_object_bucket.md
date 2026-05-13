---
subcategory : "Cloud Storage"
---

# ovh_cloud_storage_object_bucket

Creates an S3 object storage bucket in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "xxxxxxxxxx"
  name         = "my-bucket"
  region       = "GRA"

  versioning {
    status = "ENABLED"
  }

  encryption {
    algorithm = "AES256"
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `name` - (Required) Bucket name (must be globally unique and DNS-compatible).
* `region` - (Required) Region where the bucket will be created.
* `encryption` - (Optional) Server-side encryption configuration:
  * `algorithm` - (Required) Encryption algorithm (e.g. `AES256`).
* `versioning` - (Optional) Versioning configuration:
  * `status` - (Required) Versioning status (`DISABLED`, `ENABLED`, `SUSPENDED`).
* `object_lock` - (Optional) Object lock (WORM) configuration; requires versioning to be enabled:
  * `mode` - (Required) Object lock retention mode (`COMPLIANCE`, `GOVERNANCE`).
  * `retention_days` - (Required) Number of days to retain objects.
  * `retention_years` - (Optional) Number of years to retain objects (alternative to retention_days).
* `tags` - (Optional) Metadata tags for the bucket (key-value map of strings).
* `owner_user_id` - (Optional) Owner user identifier.

## Attributes Reference

The following attributes are exported:

* `id` - Bucket identifier returned by the API. In multi-region deployments this is the composite form `UPPERCASE_REGION_BUCKETNAME` (e.g. `GRA_my-bucket`); in single-region deployments it is just the bucket name.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the bucket.
* `updated_at` - Last update date of the bucket.
* `resource_status` - Bucket readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the bucket:
  * `name` - Bucket name.
  * `location` - Location details:
    * `region` - Region code.
  * `encryption` - Encryption configuration:
    * `algorithm` - Encryption algorithm.
  * `versioning` - Versioning configuration:
    * `status` - Versioning status.
  * `object_lock` - Object lock configuration:
    * `mode` - Object lock retention mode.
    * `retention_days` - Retention period in days.
    * `retention_years` - Retention period in years.
  * `tags` - Metadata tags.

## Import

A cloud storage object bucket can be imported using its `service_name`, `region` and API-side bucket `id`, separated by `/`. The API-side `id` is the composite form `UPPERCASE_REGION_BUCKETNAME` in multi-region deployments (e.g. `GRA_my-bucket`) and just the bucket name in single-region deployments:

```terraform
import {
  to = ovh_cloud_storage_object_bucket.bucket
  id = "<service_name>/<region>/<id>"
}
```

```bash
$ terraform import ovh_cloud_storage_object_bucket.bucket service_name/region/id
```
