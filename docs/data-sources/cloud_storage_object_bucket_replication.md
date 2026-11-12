---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_replication (Data Source)

Get the replication configuration of an S3&trade; compatible object storage bucket in a public cloud project.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

A bucket without replication configuration is not an error: `rules` is an empty list.

## Example Usage

```terraform
data "ovh_cloud_storage_object_bucket_replication" "replication" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-source-bucket"
}

output "replication_destinations" {
  value = [for rule in data.ovh_cloud_storage_object_bucket_replication.replication.rules : rule.destination.name]
}
```

## Argument Reference

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `bucket_id` - (Required) ID of the source bucket, in the `<REGION>_<name>` form (e.g. `GRA_my-bucket`): the `id` of an `ovh_cloud_storage_object_bucket`.

## Attributes Reference

* `id` - Replication configuration identifier, the same value as `bucket_id`.
* `rules` - List of replication rules, empty when the bucket has no replication configuration:
  * `id` - Rule identifier.
  * `status` - Rule status (`DISABLED`, `ENABLED`).
  * `priority` - Rule priority.
  * `delete_marker_replication` - Whether delete markers are replicated (`DISABLED`, `ENABLED`). Always `DISABLED` when `filter.tags` is set.
  * `destination` - Destination of the replicated objects:
    * `name` - Name of the destination bucket.
    * `region` - Region of the destination bucket, `null` when the name matches no bucket of the project, or several.
    * `storage_class` - Storage class of the replicated objects in the destination bucket (`DEEP_ARCHIVE`, `GLACIER`, `GLACIER_IR`, `HIGH_PERF`, `INTELLIGENT_TIERING`, `ONEZONE_IA`, `STANDARD`, `STANDARD_IA`).
  * `filter` - Objects the rule applies to, `null` when the rule applies to every object of the bucket:
    * `prefix` - Only objects whose key starts with this prefix.
    * `tags` - Only objects carrying every one of these tags, as a map of strings.

Optional attributes are `null` when the API does not return them.
