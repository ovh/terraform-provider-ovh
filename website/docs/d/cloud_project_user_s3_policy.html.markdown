---
layout: "ovh"
page_title: "OVH: ovh_cloud_project_user_s3_policy"
sidebar_current: "docs-ovh-datasource-cloud-project-user-s3-policy"
description: |-
  Get the S3 Policy of a public cloud project user.
---

# ovh_cloud_project_user_s3_policy

Get the S3 Policy of a public cloud project user. The policy can be set by using the `ovh_cloud_project_user_s3_policy` resource.

## Example Usage

```hcl
resource "ovh_cloud_project_user" "user" {
 service_name = "XXX"
 description  = "my user"
 role_names   = [
  "objectstore_operator"
 ]
}

resource "ovh_cloud_project_user_s3_credential" "my_s3_credentials" {
 service_name = ovh_cloud_project_user.user.service_name
 user_id      = ovh_cloud_project_user.user.id
}

resource "ovh_cloud_project_user_s3_policy" "policy" {
 service_name = ovh_cloud_project_user.user.service_name
 user_id      = ovh_cloud_project_user.user.id
 policy       = jsonencode({
  "Statement":[{
    "Sid": "RWContainer",
    "Effect": "Allow",
    "Action":["s3:GetObject", "s3:PutObject", "s3:DeleteObject", "s3:ListBucket", "s3:ListMultipartUploadParts", "s3:ListBucketMultipartUploads", "s3:AbortMultipartUpload", "s3:GetBucketLocation"],
    "Resource":["arn:aws:s3:::hp-bucket", "arn:aws:s3:::hp-bucket/*"]
  }]
 })
}

data "ovh_cloud_project_user_s3_policy" "policy" {
 service_name = ovh_cloud_project_user.user.service_name
 user_id      = ovh_cloud_project_user.user.id
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

- `user_id` - (Required) The ID of a public cloud project's user.

## Attributes Reference

- `policy` - (Required) The policy document. This is a JSON formatted string.
