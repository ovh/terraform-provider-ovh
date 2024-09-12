---
subcategory : "Cloud Project"
---

# ovh_cloud_project_user

Get the user details of a previously created public cloud project user.

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

data "ovh_cloud_project_user" "my_user" {
  service_name = data.ovh_cloud_project_users.project_users.service_name
  user_id      = local.s3_user_id
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

- `user_id` - (Required) The ID of a public cloud project's user.

## Attributes Reference

`id` is set with the user_id of the user.
In addition, the following attributes are exported:

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
