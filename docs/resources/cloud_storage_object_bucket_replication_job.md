---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_replication_job

Starts a batch replication job on an S3&trade; compatible object storage bucket of a public cloud project: the current replication rules of the bucket are applied to the objects already in it, which the rules alone do not replicate.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

The bucket must have replication rules when the job is created, and a bucket runs one job at a time: creating a job while another one is in progress fails. The job runs in the background; `terraform apply` does not wait for it to end, and every refresh updates its status and counters.

~> **Note:** A job cannot be cancelled or deleted. Destroying the resource only removes it from the Terraform state, and the job keeps running. Every argument forces a replacement, so changing one starts a new job. Once read, a job that the API no longer returns is kept in the state as is, so that its disappearance does not start a new job.

## Example Usage

```terraform
resource "ovh_cloud_storage_object_bucket_replication_job" "catch_up" {
  service_name = ovh_cloud_storage_object_bucket_replication.replication.service_name
  bucket_id    = ovh_cloud_storage_object_bucket_replication.replication.bucket_id
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value starts a new job.**
* `bucket_id` - (Required) ID of the source bucket, in the `<REGION>_<name>` form (e.g. `GRA_my-bucket`): the `id` of an `ovh_cloud_storage_object_bucket`. Reference the `ovh_cloud_storage_object_bucket_replication` resource of the bucket, as in the example, so that the rules exist before the job is created. **Changing this value starts a new job.**

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the job.
* `status` - Job status. One of `ACTIVE`, `COMPLETED`, `CREATED`, `FAILED`, `PREPARING`.
* `total` - Total number of objects to replicate.
* `completed` - Number of successfully replicated objects.
* `failed` - Number of objects that failed to replicate.
* `stuck` - Number of stuck objects.
* `created_at` - Job creation date.

## Import

A replication job can be imported using the `service_name`, the `bucket_id` and the job `id`, separated by `/`. Importing a job does not start a new one:

```terraform
import {
  to = ovh_cloud_storage_object_bucket_replication_job.catch_up
  id = "<service_name>/<bucket_id>/<job_id>"
}
```

```bash
$ terraform import ovh_cloud_storage_object_bucket_replication_job.catch_up service_name/bucket_id/job_id
```
