---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_user

Creates an user for a database cluster associated with a public cloud project.

With this resource you can create a user and map "avnadmin" for the following database engine:

  * `cassandra`
  * `kafka`
  * `kafkaConnect`
  * `mysql`
  * `grafana`

## Example Usage

Create a user johndoe in a database.
Output the user generated password with command `terraform output user_password`.

```hcl
data "ovh_cloud_project_database" "db" {
  service_name  = "XXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

resource "ovh_cloud_project_database_user" "user" {
  service_name  = data.ovh_cloud_project_database.db.service_name
  engine        = data.ovh_cloud_project_database.db.engine
  cluster_id    = data.ovh_cloud_project_database.db.id
  name          = "johndoe"
}

output "user_password" {
  value     = ovh_cloud_project_database_user.user.password
  sensitive = true
}
```

-> __NOTE__ To reset password of the user previously created, update the `password_reset` attribute.
Use the `terraform refresh` command after executing `terraform apply` to update the output with the new password.
```hcl
data "ovh_cloud_project_database" "db" {
  service_name  = "XXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

resource "ovh_cloud_project_database_user" "user" {
  service_name    = data.ovh_cloud_project_database.db.service_name
  engine          = data.ovh_cloud_project_database.db.engine
  cluster_id      = data.ovh_cloud_project_database.db.id
  name            = "johndoe"
  password_reset  = "reset1"
}

output "user_password" {
  value     = ovh_cloud_project_database_user.user.password
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required, Forces new resource) The engine of the database cluster you want to add. You can find the complete list of available engine in the [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).
Available engines:
  * `cassandra`
  * `kafka`
  * `kafkaConnect`
  * `mysql`
  * `grafana`


* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `name` - (Required, Forces new resource) Name of the user. A user named "avnadmin" is map with already created admin user and reset his password instead of create a new user. The "Grafana" engine only allows the "avnadmin" mapping.

* `password_reset` - (Optional) Arbitrary string to change to trigger a password update. Use the `terraform refresh` command after executing `terraform apply` to update the output with the new password.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `engine` - See Argument Reference above.
* `id` - ID of the user.
* `name` - See Argument Reference above.
* `password` - (Sensitive) Password of the user.
* `password_reset` - Arbitrary string to change to trigger a password update.
* `service_name` - See Argument Reference above.
* `status` - Current status of the user.

## Timeouts

```hcl
resource "ovh_cloud_project_database_user" "user" {
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

OVHcloud Managed database clusters users can be imported using the `service_name`, `engine`, `cluster_id` and `id` of the user, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_database_user.my_user service_name/engine/cluster_id/id
```
