---
subcategory: "Object Storage"
---

# ovh_cloud_s3_buckets (Data Source)

List the S3&trade; compatible object storage buckets in a public cloud project.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

## Example Usage

```terraform
data "ovh_cloud_s3_buckets" "buckets" {
  service_name = <Public cloud project id>
}
```

Filter the buckets by region:

```terraform
data "ovh_cloud_s3_buckets" "buckets" {
  service_name = <Public cloud project id>
  region       = "GRA"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `region` - (Optional) If set, only buckets located in this region are returned.

## Attributes Reference

* `buckets` - List of buckets:
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
    * `virtual_host` - Bucket virtual host. Only returned on a single bucket read, `null` here.
    * `objects_count` - Bucket total objects count. Only returned on a single bucket read, `null` here.
    * `objects_size` - Bucket total objects size in bytes. Only returned on a single bucket read, `null` here.
