---
subcategory : "Cloud Project"
---

# ovh_cloud_project_users

Get the list of all users of a public cloud project.

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

output "user_id" {
  value = local.s3_user_id
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

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
  - `username` - the username generated for the user. This username can be used with the Openstack API.
