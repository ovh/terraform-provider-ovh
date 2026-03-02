---
subcategory : "Managed Databases"
---

~> **DEPRECATED:** Use `ovh_cloud_managed_database_mongodb_prometheus` instead. This resource will be removed in the next major version.

# ovh_cloud_project_database_mongodb_prometheus

Creates a prometheus for a MongoDB cluster associated with a public cloud project.

## Example Usage

Create a Prometheus for a database. Output the prometheus user generated password with command `terraform output prom_password`.

```terraform
data "ovh_cloud_project_database" "mongodb" {
  service_name  = "XXX"
  engine        = "mongodb"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_mongodb_prometheus" "prometheus" {
  service_name  = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_project_database.mongodb.id
}

output "prom_password" {
  value     = ovh_cloud_project_database_mongodb_prometheus.prometheus.password
  sensitive = true
}
```

-> **NOTE** To reset password of the prometheus user previously created, update the `password_reset` attribute. Use the `terraform refresh` command after executing `terraform apply` to update the output with the new password. This attribute can be an arbitrary string but we recommend 2 formats:
- a datetime to keep a trace of the last reset
- a md5 of other variables to automatically trigger it based on this variable update

```terraform
data "ovh_cloud_project_database" "mongodb" {
  service_name  = "XXXX"
  engine        = "mongodb"
  id            = "ZZZZ"
}

# Change password_reset with the datetime each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_mongodb_prometheus" "prometheus_datetime" {
  service_name    = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id      = data.ovh_cloud_project_database.mongodb.id
  password_reset  = "2024-01-02T11:00:00Z"
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_project_database_mongodb_prometheus" "prometheus_md5" {
  service_name    = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id      = data.ovh_cloud_project_database.mongodb.id
  password_reset  = md5(var.something)
}

# Change password_reset each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_mongodb_prometheus" "prometheus" {
  service_name    = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id      = data.ovh_cloud_project_database.mongodb.id
  password_reset  = "reset1"
}

output "prom_password" {
  value     = ovh_cloud_project_database_mongodb_prometheus.prometheus.password
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `cluster_id` - (Required, Forces new resource) Cluster ID.
* `password_reset` - (Optional) Arbitrary string to change to trigger a password update. Use the `terraform refresh` command after executing `terraform apply` to update the output with the new password.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `engine` - See Argument Reference above.
* `id` - Cluster ID.
* `password` - (Sensitive) Password of the user.
* `password_reset` - Arbitrary string to change to trigger a password update.
* `service_name` - See Argument Reference above.
* `username` - name of the prometheus user.
* `srv_domain` - Name of the srv domain endpoint.

## Timeouts

```terraform
resource "ovh_cloud_project_database_mongodb_prometheus" "prometheus" {
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

OVHcloud Managed MongoDB clusters prometheus can be imported using the `service_name` and `cluster_id`, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_database_mongodb_prometheus.my_prometheus service_name/engine/cluster_id
```
