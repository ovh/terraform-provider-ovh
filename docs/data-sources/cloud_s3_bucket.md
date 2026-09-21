---
subcategory : "Object Storage"
---

# ovh_cloud_s3_bucket (Data Source)

Get an S3&trade; compatible object storage bucket in a public cloud project.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

## Example Usage

Look the bucket up by name. `region` is required in this form:

```hcl
data "ovh_cloud_s3_bucket" "bucket" {
  service_name = <Public cloud project id>
  name         = "my-data-bucket"
  region       = "GRA"
}
```

Or look it up by its API identifier:

```hcl
data "ovh_cloud_s3_bucket" "bucket" {
  service_name = <Public cloud project id>
  id           = "my-data-bucket"
}
```

~> __NOTE__ The bucket `id` is the bare bucket name on a single-region API instance, and `<REGION>_<name>` (for example `GRA_my-data-bucket`) on a multi-region one. Because the provider cannot tell which mode the API instance runs in, the `name` + `region` form does not build the identifier: it lists the project's buckets and reads the identifier of the one matching that name and region. The `id` form is passed through untouched and therefore works on both. Since the listing route is paginated, prefer the `id` form on projects holding a large number of buckets.

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `id` - (Optional) The identifier of the bucket. Exactly one of `id` or `name` must be set.
* `name` - (Optional) The name of the bucket. Requires `region`, and is mutually exclusive with `id`.
* `region` - (Optional) The region identifier the bucket is located in. Only valid together with `name`.

## Attributes Reference

* `id` - Bucket identifier.
* `name` - Bucket name.
* `location` - Target geographic region of the bucket:
  * `region` - Region identifier.
* `owner_user_id` - Owner user identifier.
* `encryption` - Server-side encryption configuration:
  * `algorithm` - Encryption algorithm (`AES256`, `PLAINTEXT`).
* `versioning` - Versioning configuration:
  * `status` - Versioning status (`DISABLED`, `ENABLED`, `SUSPENDED`).
* `object_lock` - Object lock (WORM) configuration:
  * `mode` - Object lock retention mode (`COMPLIANCE`, `GOVERNANCE`).
  * `retention_days` - Number of days to retain objects.
  * `retention_years` - Number of years to retain objects.
* `tags` - Metadata tags for the bucket.
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
  * `virtual_host` - Bucket virtual host, as a hostname without scheme (for example `my-data-bucket.s3.gra.io.cloud.ovh.net`).
  * `objects_count` - Bucket total objects count.
  * `objects_size` - Bucket total objects size in bytes.
