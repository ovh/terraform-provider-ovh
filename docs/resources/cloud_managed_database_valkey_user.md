---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_valkey_user

Creates an user for a Valkey cluster associated with a public cloud project.

## Example Usage

Create a user johndoe in a Valkey database. Output the user generated password with command `terraform output user_password`.

```terraform
data "ovh_cloud_managed_database" "valkey" {
  service_name  = "XXXX"
  engine        = "valkey"
  id            = "ZZZZ"
}

resource "ovh_cloud_managed_database_valkey_user" "user" {
  service_name  = data.ovh_cloud_managed_database.valkey.service_name
  cluster_id    = data.ovh_cloud_managed_database.valkey.id
  categories    = ["+@set", "+@sortedset"]
  channels      = ["*"]
  commands      = ["+get", "-set"]
  keys          = ["data", "properties"]
  name          = "johndoe"
}

output "user_password" {
  value     = ovh_cloud_managed_database_valkey_user.user.password
  sensitive = true
}
```

-> **NOTE** To reset password of the user previously created, update the `password_reset` attribute. Use the `terraform refresh` command after executing `terraform apply` to update the output with the new password. This attribute can be an arbitrary string but we recommend 2 formats:
- a datetime to keep a trace of the last reset
- a md5 of other variables to automatically trigger it based on this variable update

```terraform
data "ovh_cloud_managed_database" "valkey" {
  service_name  = "XXXX"
  engine        = "valkey"
  id            = "ZZZZ"
}

# Change password_reset with the datetime each time you want to reset the password to trigger an update
resource "ovh_cloud_managed_database_valkey_user" "user_datetime" {
  service_name    = data.ovh_cloud_managed_database.valkey.service_name
  cluster_id      = data.ovh_cloud_managed_database.valkey.id
  categories      = ["+@set", "+@sortedset"]
  channels        = ["*"]
  commands        = ["+get", "-set"]
  keys            = ["data", "properties"]
  name            = "alice"
  password_reset  = "2024-01-02T11:00:00Z"
}

variable "something" {
  type = string
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_managed_database_valkey_user" "user_md5" {
  service_name    = data.ovh_cloud_managed_database.valkey.service_name
  cluster_id      = data.ovh_cloud_managed_database.valkey.id
  categories      = ["+@set", "+@sortedset"]
  channels        = ["*"]
  commands        = ["+get", "-set"]
  keys            = ["data", "properties"]
  name            = "bob"
  password_reset  = md5(var.something)
}

# Change password_reset each time you want to reset the password to trigger an update
resource "ovh_cloud_managed_database_valkey_user" "user" {
  service_name    = data.ovh_cloud_managed_database.valkey.service_name
  cluster_id      = data.ovh_cloud_managed_database.valkey.id
  categories      = ["+@set", "+@sortedset"]
  channels        = ["*"]
  commands        = ["+get", "-set"]
  keys            = ["data", "properties"]
  name            = "johndoe"
  password_reset  = "reset1"
}

output "user_password" {
  value = ovh_cloud_managed_database_valkey_user.user.password
  sensitive = true
}
```

To reset the password of the existing default user:
```
data "ovh_cloud_managed_database" "valkey" {
  service_name  = "XXXX"
  engine        = "valkey"
  id            = "ZZZZ"
}

resource "ovh_cloud_managed_database_valkey_user" "default_user" {
  service_name    = data.ovh_cloud_managed_database.valkey.service_name
  cluster_id      = data.ovh_cloud_managed_database.valkey.id
  categories      = ["+@all"]
  keys            = ["*"]
  name            = "default"
  password_reset  = "reset"
}

output "default_user_password" {
  value = ovh_cloud_managed_database_valkey_user.default_user.password
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `cluster_id` - (Required, Forces new resource) Cluster ID.
* `categories` - (Optional) Categories of the user.
* `channels` - (Optional: if omit, all channels) Channels of the user.
* `commands` - (Optional) Commands of the user.
* `keys` - (Optional) Keys of the user.
* `name` - (Required, Forces new resource) Name of the user. A user named "default" is mapped with an already created user and reset his password instead of creating a new user.
* `password_reset` - (Optional) Arbitrary string to change to trigger a password update. Use the `terraform refresh` command after executing `terraform apply` to update the output with the new password.

## Attributes Reference

The following attributes are exported:

* `categories` - See Argument Reference above.
* `channels` - See Argument Reference above.
* `cluster_id` - See Argument Reference above.
* `commands` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `keys` - See Argument Reference above.
* `name` - See Argument Reference above.
* `password` - (Sensitive) Password of the user.
* `password_reset` - Arbitrary string to change to trigger a password update.
* `service_name` - See Argument Reference above.
* `status` - Current status of the user.

## Timeouts

```terraform
resource "ovh_cloud_managed_database_valkey_user" "user" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}
```
* `create` - (Default 20m)
* `update` - (Default 20m)
* `delete` - (Default 20m)

## Import

OVHcloud Managed Valkey clusters users can be imported using the `service_name`, `cluster_id` and `id` of the user, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_managed_database_valkey_user.my_user service_name/cluster_id/id
```
