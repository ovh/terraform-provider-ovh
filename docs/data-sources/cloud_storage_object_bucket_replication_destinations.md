---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_replication_destinations (Data Source)

List the other buckets of a public cloud project and whether each one can receive the objects replicated from an S3&trade; compatible object storage bucket. Every page of the listing is read.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

## Example Usage

```terraform
data "ovh_cloud_storage_object_bucket_replication_destinations" "destinations" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-source-bucket"
}

output "eligible_destinations" {
  value = [
    for destination in data.ovh_cloud_storage_object_bucket_replication_destinations.destinations.destinations :
    "${destination.region}/${destination.name}" if destination.eligible
  ]
}
```

## Argument Reference

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `bucket_id` - (Required) ID of the source bucket, in the `<REGION>_<name>` form (e.g. `GRA_my-bucket`): the `id` of an `ovh_cloud_storage_object_bucket`.

## Attributes Reference

* `destinations` - List of the buckets of the project evaluated as destinations of the source bucket:
  * `name` - Bucket name.
  * `region` - Region the bucket is located in.
  * `eligible` - Whether the bucket can be the destination of a replication rule on the source bucket.
  * `reason` - Why the bucket cannot be a destination, `null` when `eligible` is `true`. One of `OBJECT_LOCK_INCOMPATIBLE` (the object lock configuration of the bucket is incompatible with the one of the source bucket), `REGION_NOT_REPLICATION_TARGET` (the region of the bucket does not accept replicated objects), `VERSIONING_DISABLED` (versioning is not enabled on the bucket).
