---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_postgresql_connection_pool (Data Source)

Use this data source to get information about a connection pool of a postgresql cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_managed_database_postgresql_connection_pool" "test_pool" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "test_pool" {
  value = {
    service_name: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.service_name
    cluster_id: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.cluster_id
    name: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.name
    database_id: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.database_id
    mode: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.mode
    size: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.size
    port: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.port
    ssl_mode: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.ssl_mode
    uri: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.uri
    user_id: data.ovh_cloud_managed_database_postgresql_connection_pool.test_pool.user_id
  }
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `cluster_id` - (Required) Cluster ID.
* `name` - (Required) Name of the Connection pool.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above
* `cluster_id` - See Argument Reference above
* `name` - See Argument Reference above
* `database_id` - Database ID for a database that belongs to the Database cluster given above.
* `mode` - Connection mode to the connection pool Available modes:
  * `session`
  * `statement`
  * `transaction`
* `size` - Size of the connection pool.
* `user_id` - Database user authorized to connect to the pool, if none all the users are allowed.
* `port` - Port of the connection pool.
* `ssl_mode` - Ssl connection mode for the pool.
* `uri` - Connection URI to the pool.
