---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_database

Creates a database for a database cluster associated with a public cloud project.

With this resource you can create a database for the following database engine:

  * `mysql`
  * `postgresql`

## Example Usage

```hcl
data "ovh_cloud_project_database" "db" {
  service_name  = "XXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

resource "ovh_cloud_project_database_database" "database" {
  service_name  = data.ovh_cloud_project_database.db.service_name
  engine        = data.ovh_cloud_project_database.db.engine
  cluster_id    = data.ovh_cloud_project_database.db.id
  name          = "mydatabase"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required, Forces new resource) The engine of the database cluster you want to add. You can find the complete list of available engine in the [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).
Available engines:
  * `mysql`
  * `postgresql`

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `name` - (Required, Forces new resource) Name of the database.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `default` - Defines if the database has been created by default.
* `engine` - See Argument Reference above.
* `id` - ID of the database.
* `service_name` - See Argument Reference above.
* `name` - See Argument Reference above.

## Timeouts

```hcl
resource "ovh_cloud_project_database_database" "database" {
  # ...

  timeouts {
    create = "1h"
    delete = "45m"
  }
}
```
* `create` - (Default 20m)
* `delete` - (Default 20m)

## Import

OVHcloud Managed database clusters databases can be imported using the `service_name`, `engine`, `cluster_id` and `id` of the database, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_database_database.my_database service_name/engine/cluster_id/id
```