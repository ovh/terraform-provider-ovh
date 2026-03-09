---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_postgresql_connection_pool

Creates a connection pool for a PostgreSQL cluster associated with a public cloud project.

## Example Usage

Create a `johndoe` user in a `mydatabase` database located in a `test-postgresql-cluster` PostgreSQL 15 cluster. Then creates a `test_connection_pool` connection pool for this user to connect to the DB. Output the generated connection pool params with command `terraform output test_pool`.

```terraform
resource "ovh_cloud_managed_database" "db" {
  service_name  = "XXXX"
  engine        = "postgresql"
  description  = "test-postgresql-cluster"
  version      = "15"
  plan         = "essential"
  nodes {
    region     = "GRA"
  }
  flavor = "db1-4"
}

resource "ovh_cloud_managed_database_database" "database" {
  service_name  = ovh_cloud_managed_database.db.service_name
  engine        = ovh_cloud_managed_database.db.engine
  cluster_id    = ovh_cloud_managed_database.db.id
  name          = "mydatabase"
}

resource "ovh_cloud_managed_database_postgresql_user" "user" {
  service_name = ovh_cloud_managed_database.db.service_name
  cluster_id   = ovh_cloud_managed_database.db.id
  name          = "johndoe"
  roles         = ["replication"]
}

resource "ovh_cloud_managed_database_postgresql_connection_pool" "test_pool" {
  service_name = ovh_cloud_managed_database.db.service_name
  cluster_id   = ovh_cloud_managed_database.db.id
  database_id = ovh_cloud_managed_database_database.database.id
  name = "test_connection_pool"
  user_id = ovh_cloud_managed_database_postgresql_user.user.id
  mode = "session"
  size = 13
}

output "test_pool" {
  value = {
    service_name: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.service_name
    cluster_id: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.cluster_id
    name: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.name
    database_id: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.database_id
    mode: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.mode
    size: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.size
    port: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.port
    ssl_mode: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.ssl_mode
    uri: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.uri
    user_id: ovh_cloud_managed_database_postgresql_connection_pool.test_pool.user_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `cluster_id` - (Required, Forces new resource) Cluster ID.
* `database_id` - (Required, Forces new resource) Database ID for a database that belongs to the Database cluster given above.
* `name` - (Required, Forces new resource) Name of the connection pool.
* `mode` - (Required) Connection mode to the connection pool Available modes:
  * `session`
  * `statement`
  * `transaction`
* `size` - (Required) Size of the connection pool.
* `user_id` - (Optional) Database user authorized to connect to the pool, if none all the users are allowed.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `cluster_id` - See Argument Reference above.
* `database_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `mode` - See Argument Reference above.
* `size` - See Argument Reference above.
* `user_id` - See Argument Reference above.
* `port` - Port of the connection pool.
* `ssl_mode` - Ssl connection mode for the pool.
* `uri` - Connection URI to the pool.

## Timeouts

```terraform
resource "ovh_cloud_managed_database_postgresql_connection_pool" "user" {
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

OVHcloud Managed PostgreSQL clusters connection pools can be imported using the `service_name`, `cluster_id` and `id` of the connection pool, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_managed_database_postgresql_connection_pool.my_connection_pool service_name/cluster_id/id
```
