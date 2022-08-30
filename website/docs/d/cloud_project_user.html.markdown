---
layout: "ovh"
page_title: "OVH: ovh_cloud_project_user"
sidebar_current: "docs-ovh-datasource-cloud-project-user"
description: |-
  Get the details of a public cloud project user.
---

# ovh_cloud_project_user

Get the user details of a previously created public cloud project user.

## Example Usage

```hcl
resource "ovh_cloud_project_user" "user" {
 service_name = "XXX"
 description  = "my user"
 role_names   = [
  "objectstore_operator"
 ]
}

data "ovh_cloud_project_user" "my_user" {
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
