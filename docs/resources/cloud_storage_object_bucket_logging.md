---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_logging

Manages the server access logging of an S3&trade; compatible object storage bucket in a public cloud project: the delivery of a record of every request made to the bucket into another bucket of the same region.

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

OVHcloud only exposes server access logging through the S3 API, so this resource needs the S3 credentials of a user of the project in the provider configuration (`s3_access_key_id` and `s3_secret_access_key`, or the `OVH_S3_ACCESS_KEY_ID` and `OVH_S3_SECRET_ACCESS_KEY` environment variables). That user must be allowed to read and write the logging configuration of the bucket (`s3:GetBucketLogging` and `s3:PutBucketLogging`), which the owner of the bucket is.

The target bucket must grant the log delivery group the right to write into it: the `log-delivery-write` canned ACL, set with `ovh_cloud_storage_object_bucket_acl`. The first logs arrive about an hour after the logging is enabled. Destroying this resource disables the logging of the bucket. If the logging is disabled outside of Terraform, the next `terraform apply` enables it again. A refresh during a region maintenance or an object storage outage keeps the last known state and shows a warning.

## Example Usage

```terraform
resource "ovh_cloud_storage_object_bucket" "logs" {
  service_name = "<Public cloud project id>"
  name         = "my-logs-bucket"
  region       = "GRA"
}

resource "ovh_cloud_storage_object_bucket_acl" "logs" {
  service_name = ovh_cloud_storage_object_bucket.logs.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.logs.id
  acl          = "log-delivery-write"
}

resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "<Public cloud project id>"
  name         = "my-data-bucket"
  region       = "GRA"
}

resource "ovh_cloud_storage_object_bucket_logging" "logging" {
  service_name  = ovh_cloud_storage_object_bucket.bucket.service_name
  bucket_id     = ovh_cloud_storage_object_bucket.bucket.id
  target_bucket = ovh_cloud_storage_object_bucket.logs.name
  target_prefix = "my-data-bucket/"

  depends_on = [ovh_cloud_storage_object_bucket_acl.logs]
}
```

## Argument Reference

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `bucket_id` - (Required) The bucket identifier, in the `<REGION>_<name>` form (the `id` of an `ovh_cloud_storage_object_bucket`). Changing it creates a new resource.
* `target_bucket` - (Required) The name of the bucket receiving the logs. It must be in the same region as the bucket, with the `log-delivery-write` canned ACL.
* `target_prefix` - (Optional) The prefix of the keys of the log objects. Defaults to `""`.

## Attributes Reference

* `id` - The identifier of the logging configuration, the bucket identifier.

## Import

The logging configuration of a bucket can be imported using the `service_name` and the `bucket_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_storage_object_bucket_logging.logging
  id = "<service_name>/<bucket_id>"
}
```

```bash
$ terraform import ovh_cloud_storage_object_bucket_logging.logging service_name/bucket_id
```
