---
subcategory : "Managed Databases"
---

~> **DEPRECATED:** Use `ovh_cloud_managed_database_mongodb_user` instead. This resource will be removed in the next major version.

# ovh_cloud_project_database_mongodb_user

Creates an user for a MongoDB cluster associated with a public cloud project.

## Example Usage

Create a user johndoe in a MongoDB database. Output the user generated password with command `terraform output user_password`.

```terraform
data "ovh_cloud_project_database" "mongodb" {
  service_name  = "XXX"
  engine        = "mongodb"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_mongodb_user" "user" {
  service_name  = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_project_database.mongodb.id
  name          = "johndoe"
  roles         = ["backup@admin", "readAnyDatabase@admin"]
}

output "user_password" {
  value     = ovh_cloud_project_database_mongodb_user.user.password
  sensitive = true
}
```

-> **NOTE** To reset password of the user previously created, update the `password_reset` attribute. Use the `terraform refresh` command after executing `terraform apply` to update the output with the new password. This attribute can be an arbitrary string but we recommend 2 formats:
- a datetime to keep a trace of the last reset
- a md5 of other variables to automatically trigger it based on this variable update

```terraform
data "ovh_cloud_project_database" "mongodb" {
  service_name  = "XXX"
  engine        = "mongodb"
  id            = "ZZZ"
}

# Change password_reset with the datetime each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_mongodb_user" "user_datetime" {
  service_name  = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_project_database.mongodb.id
  name          = "alice"
  roles         = ["backup@admin", "readAnyDatabase@admin"]
  password_reset  = "2024-01-02T11:00:00Z"
}

variable "something" {
  type = string
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_project_database_mongodb_user" "user_md5" {
  service_name  = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_project_database.mongodb.id
  name          = "bob"
  roles         = ["backup@admin", "readAnyDatabase@admin"]
  password_reset  = md5(var.something)
}

# Change password_reset each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_mongodb_user" "user" {
  service_name  = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_project_database.mongodb.id
  name          = "johndoe"
  roles         = ["backup@admin", "readAnyDatabase@admin"]
  password_reset  = "reset1"
}

output "user_password" {
  value     = ovh_cloud_project_database_mongodb_user.user.password
  sensitive = true
}
```

To map the "admin@admin" user and keep default roles

```terraform
data "ovh_cloud_project_database" "mongodb" {
  service_name  = "XXX"
  engine        = "mongodb"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_mongodb_user" "user" {
  service_name  = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_project_database.mongodb.id
  name          = "admin"
  roles         = ["clusterMonitor@admin", "readWriteAnyDatabase@admin","userAdminAnyDatabase@admin"]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `cluster_id` - (Required, Forces new resource) Cluster ID.
* `name` - (Required, Forces new resource) Name of the user. A user named "admin" is mapped with already created admin@admin user instead of creating a new user.
* `roles` - (Optional: if omit, default role) Roles the user belongs to. Since version 0.37.0, the authentication database must be indicated for all roles. Available roles:
  * `backup@admin`
  * `clusterAdmin@admin`
  * `clusterManager@admin`
  * `clusterMonitor@admin`
  * `dbAdmin@(defined db)`
  * `dbAdminAnyDatabase@admin`
  * `dbOwner@(defined db)`
  * `enableSharding@(defined db)`
  * `hostManager@admin`
  * `read@(defined db)`
  * `readAnyDatabase@admin`
  * `readWrite@(defined db)`
  * `readWriteAnyDatabase@admin`
  * `restore@admin`
  * `root@admin`
  * `userAdmin@(defined db)`
  * `userAdminAnyDatabase@admin`
* `password_reset` - (Optional) Arbitrary string to change to trigger a password update. Use the `terraform refresh` command after executing `terraform apply` to update the output with the new password.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `name` - Name of the user with the authentication database in the format name@authDB
* `password` - (Sensitive) Password of the user.
* `password_reset` - See Argument Reference above.
* `roles` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `status` - Current status of the user.

## Timeouts

```terraform
resource "ovh_cloud_project_database_mongodb_user" "user" {
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

OVHcloud Managed MongoDB clusters users can be imported using the `service_name`, `cluster_id` and `id` of the user, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_database_mongodb_user.my_user service_name/cluster_id/id
```
