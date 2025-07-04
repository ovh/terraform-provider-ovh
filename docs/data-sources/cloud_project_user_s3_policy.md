---
subcategory : "Object Storage"
---

# ovh_cloud_project_user_s3_policy

Get the S3 Policy of a public cloud project user. The policy can be set by using the `ovh_cloud_project_user_s3_policy` resource.

## Example Usage

```terraform
data "ovh_cloud_project_users" "project_users" {
  service_name = "XXX"
}

locals {
  # Get the user ID of a previously created user with the description "S3-User"
  users      = [for user in data.ovh_cloud_project_users.project_users.users : user.user_id if user.description == "S3-User"]
  s3_user_id = local.users[0]
}

data "ovh_cloud_project_user_s3_policy" "policy" {
  service_name = data.ovh_cloud_project_users.project_users.service_name
  user_id      = local.s3_user_id
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

- `user_id` - (Required) The ID of a public cloud project's user.

## Attributes Reference

- `policy` - (Required) The policy document. This is a JSON formatted string.
