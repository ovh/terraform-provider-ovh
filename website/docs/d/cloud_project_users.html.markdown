---
layout: "ovh"
page_title: "OVH: ovh_cloud_project_users"
sidebar_current: "docs-ovh-datasource-cloud-project-users"
description: |-
  Get the list of all users of a public cloud project.
---

# ovh_cloud_project_users

Get the list of all users of a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_users" "project_users" {
 service_name = "XXX"
}

locals {
  # Get the user ID of a previously created user with the description "S3-User"
  users = [for user in data.ovh_cloud_project_users.project_users.users : user.user_id if user.description == "S3-User"]
  s3_user_id = users[0]
}

resource "ovh_cloud_project_user_s3_credential" "my_s3_credentials" {
 service_name = data.ovh_cloud_project_users.project_users.service_name
 user_id      = local.s3_user_id
}

resource "ovh_cloud_project_user_s3_policy" "policy" {
 service_name = data.ovh_cloud_project_users.project_users.service_name
 user_id      = local.s3_user_id
 policy       = jsonencode({
  "Statement":[{
    "Sid": "RWContainer",
    "Effect": "Allow",
    "Action":["s3:GetObject", "s3:PutObject", "s3:DeleteObject", "s3:ListBucket", "s3:ListMultipartUploadParts", "s3:ListBucketMultipartUploads", "s3:AbortMultipartUpload", "s3:GetBucketLocation"],
    "Resource":["arn:aws:s3:::hp-bucket", "arn:aws:s3:::hp-bucket/*"]
  }]
 })
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

## Attributes Reference

- `users` - The list of users of a public cloud project.
  - `user_id` - The ID of a public cloud project's user.
  - `creation_date` - the date the user was created.
  - `description` - See Argument Reference above.
  - `roles` - A list of roles associated with the user.
    - `description` - description of the role
    - `id` - id of the role
    - `name` - name of the role
    - `permissions` - list of permissions associated with the role
  - `status` - the status of the user. should be normally set to 'ok'.
  - `username` - the username generated for the user. This username can be used with
    the Openstack API.
