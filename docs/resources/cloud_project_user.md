---
subcategory : "Cloud Project"
---

# ovh_cloud_project_user

Creates a user in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_project_user" "user1" {
  service_name = "XXX"
}
```

## Argument Reference

The following arguments are supported:

* `description` - A description associated with the user.

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `role_name` - The name of a role. See `role_names`.

* `role_names` - A list of role names. Values can be:
  - administrator,
  - ai_training_operator
  - ai_training_read
  - authentication
  - backup_operator
  - compute_operator
  - image_operator
  - infrastructure_supervisor
  - network_operator
  - network_security_operator
  - objectstore_operator
  - volume_operator

* `rotate_when_changed` - A map of arbitrary key/value pairs that will trigger password regeneration when they change, enabling password rotation based on external conditions such as a rotating timestamp. When rotation occurs, the new password is saved in the Terraform state and can be used by dependent resources (such as storing the password in a key vault).

## Attributes Reference

The following attributes are exported:

* `creation_date` - the date the user was created.
* `description` - See Argument Reference above.
* `openstack_rc` - a convenient map representing an openstack_rc file. Note: no password nor sensitive token is set in this map.
* `password` - (Sensitive) the password generated for the user. The password can be used with the Openstack API. This attribute is sensitive and will only be retrieve once during creation.
* `roles` - A list of roles associated with the user.
  * `description` - description of the role
  * `id` - id of the role
  * `name` - name of the role
  * `permissions` - list of permissions associated with the role
* `service_name` - See Argument Reference above.
* `status` - the status of the user. should be normally set to 'ok'.
* `username` - the username generated for the user. This username can be used with the Openstack API.
