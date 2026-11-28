---
subcategory : "Object Storage"
---

# ovh_cloud_storage_object_bucket_acl

Manages the canned ACL of an S3&trade; compatible object storage bucket in a public cloud project: `private`, `public-read`, `public-read-write`, `authenticated-read` or `log-delivery-write`, the ACL a bucket needs to receive the server access logs of other buckets (`ovh_cloud_storage_object_bucket_logging`).

~> __NOTE__ S3 is a trademark filed by Amazon Technologies, Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies, Inc.

OVHcloud only exposes bucket ACLs through the S3 API, so this resource needs the S3 credentials of a user of the project in the provider configuration (`s3_access_key_id` and `s3_secret_access_key`, or the `OVH_S3_ACCESS_KEY_ID` and `OVH_S3_SECRET_ACCESS_KEY` environment variables). That user must be allowed to read and write the ACL of the bucket (`s3:GetBucketAcl` and `s3:PutBucketAcl`), which the owner of the bucket is.

Destroying this resource sets the bucket back to `private`, the ACL of a new bucket. If the grants of the bucket are changed outside of Terraform to another canned ACL, the next `terraform apply` sets it back. Grants that amount to no canned ACL cannot be described here: the refresh keeps the last known state and shows a warning. A refresh during a region maintenance or an object storage outage also keeps the last known state and shows a warning.

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
```

## Argument Reference

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `bucket_id` - (Required) The bucket identifier, in the `<REGION>_<name>` form (the `id` of an `ovh_cloud_storage_object_bucket`). Changing it creates a new resource.
* `acl` - (Required) The canned ACL of the bucket: `authenticated-read`, `log-delivery-write`, `private`, `public-read` or `public-read-write`.

## Attributes Reference

* `id` - The identifier of the ACL, the bucket identifier.

## Import

The ACL of a bucket can be imported using the `service_name` and the `bucket_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_storage_object_bucket_acl.acl
  id = "<service_name>/<bucket_id>"
}
```

```bash
$ terraform import ovh_cloud_storage_object_bucket_acl.acl service_name/bucket_id
```
