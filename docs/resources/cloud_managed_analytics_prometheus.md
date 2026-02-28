---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_prometheus

Creates a prometheus for a database cluster associated with a public cloud project.

With this resource you can create a prometheus for the following database engine:

* `cassandra`
* `kafka`
* `kafkaConnect`
* `kafkaMirrorMaker`
* `mysql`
* `opensearch`
* `postgresql`
* `redis`
* `valkey`

## Example Usage

Create a Prometheus for a database. Output the prometheus user generated password with command `terraform output prom_password`.

```terraform
data "ovh_cloud_managed_analytics" "db" {
  service_name  = "XXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

resource "ovh_cloud_managed_analytics_prometheus" "prometheus" {
  service_name  = data.ovh_cloud_managed_analytics.db.service_name
  engine        = data.ovh_cloud_managed_analytics.db.engine
  cluster_id    = data.ovh_cloud_managed_analytics.db.id
}

output "prom_password" {
  value     = ovh_cloud_managed_analytics_prometheus.prometheus.password
  sensitive = true
}
```

-> **NOTE** To reset password of the prometheus user previously created, update the `password_reset` attribute. Use the `terraform refresh` command after executing `terraform apply` to update the output with the new password. This attribute can be an arbitrary string but we recommend 2 formats:
- a datetime to keep a trace of the last reset
- a md5 of other variables to automatically trigger it based on this variable update

```terraform
data "ovh_cloud_managed_analytics" "db" {
  service_name  = "XXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_managed_analytics_prometheus" "prometheus_datetime" {
  service_name    = data.ovh_cloud_managed_analytics.db.service_name
  engine          = data.ovh_cloud_managed_analytics.db.engine
  cluster_id      = data.ovh_cloud_managed_analytics.db.id
  password_reset  = "2024-01-02T11:00:00Z"
}

variable "something" {
  type = string
}

resource "ovh_cloud_managed_analytics_prometheus" "prometheus_md5" {
  service_name    = data.ovh_cloud_managed_analytics.db.service_name
  engine          = data.ovh_cloud_managed_analytics.db.engine
  cluster_id      = data.ovh_cloud_managed_analytics.db.id
  password_reset  = md5(var.something)
}

resource "ovh_cloud_managed_analytics_prometheus" "prometheus" {
  service_name    = data.ovh_cloud_managed_analytics.db.service_name
  engine          = data.ovh_cloud_managed_analytics.db.engine
  cluster_id      = data.ovh_cloud_managed_analytics.db.id
  password_reset  = "reset1"
}

output "prom_password" {
  value     = ovh_cloud_managed_analytics_prometheus.prometheus.password
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `engine` - (Required, Forces new resource) The engine of the database cluster you want to add. You can find the complete list of available engine in the [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases). Available engines:
  * `cassandra`
  * `kafka`
  * `kafkaConnect`
  * `kafkaMirrorMaker`
  * `mysql`
  * `opensearch`
  * `postgresql`
  * `redis`
  * `valkey`
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
* `targets` - List of all endpoint targets.
  * `Host` - Host of the endpoint.
  * `Port` - Connection port for the endpoint.
* `username` - name of the prometheus user.

## Timeouts

```terraform
resource "ovh_cloud_managed_analytics_prometheus" "prometheus" {
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

OVHcloud Managed database clusters prometheus can be imported using the `service_name`, `engine` and `cluster_id`, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_managed_analytics_prometheus.my_prometheus service_name/engine/cluster_id
```
