---
subcategory : "Account Management"
---

# ovh_cloud_project_user_s3_credential (Data Source)

Use this data source to retrieve the Secret Access Key of an Access Key ID associated with a public cloud project's user.

## Example Usage

```hcl
data "ovh_cloud_project_users" "project_users" {
  service_name = "XXX"
}

locals {
  # Get the user ID of a previously created user with the description "S3-User"
  users      = [for user in data.ovh_cloud_project_users.project_users.users : user.user_id if user.description == "S3-User"]
  s3_user_id = local.users[0]
}

data "ovh_cloud_project_user_s3_credentials" "my_s3_credentials" {
  service_name = data.ovh_cloud_project_users.project_users.service_name
  user_id      = local.s3_user_id
}

data "ovh_cloud_project_user_s3_credential" "my_s3_credential" {
  service_name  = data.ovh_cloud_project_user_s3_credentials.my_s3_credentials.service_name
  user_id       = data.ovh_cloud_project_user_s3_credentials.my_s3_credentials.user_id
  access_key_id = data.ovh_cloud_project_user_s3_credentials.my_s3_credentials.access_key_ids[0]
}

output "my_access_key_id" {
  value = data.ovh_cloud_project_user_s3_credential.my_s3_credential.access_key_id
}

output "my_secret_access_key" {
  value     = data.ovh_cloud_project_user_s3_credential.my_s3_credential.secret_access_key
  sensitive = true
}
```

## Argument Reference

- `service_name` - (Required) The ID of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

- `user_id` - (Required) The ID of a public cloud project's user.

- `access_key_id` - the Access Key ID

## Attributes Reference

- `secret_access_key` - (Sensitive) the Secret Access Key
