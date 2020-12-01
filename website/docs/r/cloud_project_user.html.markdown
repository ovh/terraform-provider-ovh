---
layout: "ovh"
page_title: "OVH: ovh_cloud_project_user"
sidebar_current: "docs-ovh-resource-cloud-project-user"
description: |-
  Creates a user in a public cloud project.
---

# ovh_cloud_project_user

Creates a user in a public cloud project.

## Example Usage

```hcl
resource "ovh_cloud_project_user" "user1" {
   project_id = "67890"
}
```

## Argument Reference

The following arguments are supported:

* `description` - A description associated with the user.

* `project_id` - (Optional) Deprecated. The id of the public cloud project. If omitted,
    the `OVH_PROJECT_ID` environment variable is used.
    One of `service_name` or `project_id` is required. Conflits with `service_name`.

* `service_name` - (Optional) The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. 
    One of `service_name` or `project_id` is required. Conflits with `project_id`.

* `role_name` -  The name of a role. See `role_names`.

* `role_names` - A list of role names. Values can be: 
  - administrator,
  - ai_training_operator
  - authentication
  - backup_operator
  - compute_operator
  - image_operator 
  - infrastructure_supervisor
  - network_operator
  - network_security_operator
  - objectstore_operator
  - volume_operator

* `service_name` -  The id of the public cloud project. Conflicts with `project_id`.


## Attributes Reference

The following attributes are exported:

* `creation_date` - the date the user was created.
* `description` - See Argument Reference above.
* `openstack_rc` - a convenient map representing an openstack_rc file.
   Note: no password nor sensitive token is set in this map.
* `password` - (Sensitive) the password generated for the user. The password can
   be used with the Openstack API. This attribute is sensitive and will only be
   retrieve once during creation.
* `project_id` - See Argument Reference above.
* `roles` - A list of roles associated with the user.
  * `description` - description of the role
  * `id` - id of the role
  * `name` - name of the role
  * `permissions` - list of permissions associated with the role
* `service_name` - See Argument Reference above.
* `status` - the status of the user. should be normally set to 'ok'.
* `username` - the username generated for the user. This username can be used with
   the Openstack API.
