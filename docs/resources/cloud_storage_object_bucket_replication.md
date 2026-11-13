---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_replication

Manages the replication configuration of an S3&trade; compatible object storage bucket in a public cloud project: rules that copy the new objects of the bucket to other buckets of the project.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

This resource owns the whole replication configuration of the bucket. Creating and updating it replace every rule of the bucket with the `rules` list, and the change is applied synchronously. Destroying it removes every replication rule of the bucket. If the replication configuration is removed outside of Terraform, the next `terraform apply` recreates it. A refresh during a region maintenance or an object storage outage keeps the last known state and shows a warning.

Replication requires versioning on the source and the destination buckets. The `ovh_cloud_storage_object_bucket_replication_destinations` data source tells which buckets of the project can be a destination. Rules only apply to the objects written after they are set; use `ovh_cloud_storage_object_bucket_replication_job` to replicate the objects already in the bucket.

## Example Usage

```terraform
resource "ovh_cloud_storage_object_bucket" "replica" {
  service_name = "<Public cloud project id>"
  name         = "my-replica-bucket"
  region       = "SBG"

  versioning = {
    status = "ENABLED"
  }
}

resource "ovh_cloud_storage_object_bucket" "source" {
  service_name = "<Public cloud project id>"
  name         = "my-source-bucket"
  region       = "GRA"

  versioning = {
    status = "ENABLED"
  }
}

resource "ovh_cloud_storage_object_bucket_replication" "replication" {
  service_name = ovh_cloud_storage_object_bucket.source.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.source.id

  rules = [
    {
      id                        = "replicate-all"
      status                    = "ENABLED"
      priority                  = 1
      delete_marker_replication = "DISABLED"
      destination = {
        name   = ovh_cloud_storage_object_bucket.replica.name
        region = ovh_cloud_storage_object_bucket.replica.region
      }
    },
    {
      id                        = "archive-logs"
      status                    = "ENABLED"
      priority                  = 2
      delete_marker_replication = "DISABLED"
      filter = {
        prefix = "logs/"
        tags = {
          retention = "long"
        }
      }
      destination = {
        name          = ovh_cloud_storage_object_bucket.replica.name
        region        = ovh_cloud_storage_object_bucket.replica.region
        storage_class = "STANDARD_IA"
      }
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `bucket_id` - (Required) ID of the source bucket, in the `<REGION>_<name>` form (e.g. `GRA_my-bucket`): the `id` of an `ovh_cloud_storage_object_bucket`. **Changing this value recreates the resource.**
* `rules` - (Required) List of replication rules, at least one. Every change is applied in place:
  * `id` - (Required) Rule identifier, unique in the bucket.
  * `status` - (Required) Rule status. One of `DISABLED`, `ENABLED`.
  * `priority` - (Required) Rule priority.
  * `delete_marker_replication` - (Required) Whether delete markers are replicated. One of `DISABLED`, `ENABLED`. Must be `DISABLED` when `filter.tags` is set.
  * `destination` - (Required) Destination of the replicated objects:
    * `name` - (Required) Name of the destination bucket.
    * `region` - (Required) Region of the destination bucket (e.g. `GRA`), matched case-insensitively: it must be the region where the destination bucket lives.
    * `storage_class` - (Optional) Storage class of the replicated objects in the destination bucket. One of `DEEP_ARCHIVE`, `GLACIER`, `GLACIER_IR`, `HIGH_PERF`, `INTELLIGENT_TIERING`, `ONEZONE_IA`, `STANDARD`, `STANDARD_IA`. When omitted, the value reported by the API is kept in the state.
  * `filter` - (Optional) Objects the rule applies to. When omitted, the rule applies to every object of the bucket:
    * `prefix` - (Optional) Only objects whose key starts with this prefix.
    * `tags` - (Optional) Only objects carrying every one of these tags, as a map of strings.

## Attributes Reference

The following attributes are exported:

* `id` - Replication configuration identifier, the same value as `bucket_id`.

## Import

A bucket replication configuration can be imported using the `service_name` and the `bucket_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_storage_object_bucket_replication.replication
  id = "<service_name>/<bucket_id>"
}
```

```bash
$ terraform import ovh_cloud_storage_object_bucket_replication.replication service_name/bucket_id
```
