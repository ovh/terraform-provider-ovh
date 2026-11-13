---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_replication_job (Data Source)

Get a batch replication job of an S3&trade; compatible object storage bucket in a public cloud project, and its progress.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

## Example Usage

```terraform
data "ovh_cloud_storage_object_bucket_replication_job" "job" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-source-bucket"
  id           = "<job id>"
}

output "replication_job_progress" {
  value = "${data.ovh_cloud_storage_object_bucket_replication_job.job.completed}/${data.ovh_cloud_storage_object_bucket_replication_job.job.total}"
}
```

## Argument Reference

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `bucket_id` - (Required) ID of the source bucket, in the `<REGION>_<name>` form (e.g. `GRA_my-bucket`): the `id` of an `ovh_cloud_storage_object_bucket`.
* `id` - (Required) Unique identifier of the job.

## Attributes Reference

* `id` - Unique identifier of the job.
* `status` - Job status (`ACTIVE`, `COMPLETED`, `CREATED`, `FAILED`, `PREPARING`).
* `total` - Total number of objects to replicate.
* `completed` - Number of successfully replicated objects.
* `failed` - Number of objects that failed to replicate.
* `stuck` - Number of stuck objects.
* `created_at` - Job creation date.
