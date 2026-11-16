---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_postgresql_connection_pool (Data Source)

Use this data source to get a connection pool of a managed database service (postgresql), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_database_postgresql_connection_pool" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  id           = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
}
```

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `service_id` - (Required) ID of the service.
* `id` - (Required) ID of the connection pool.

## Attributes Reference

The following attributes are exported:

* `engine` - Engine of the service: postgresql.
* `name` - Name of the connection pool.
* `checksum` - Hash of the connection pool, computed by the API.
* `created_at` - Creation date of the connection pool.
* `database_id` - ID of the database the pool connects to.
* `mode` - PgBouncer pooling mode: SESSION, STATEMENT or TRANSACTION.
* `port` - Port the connection pool listens on.
* `size` - Maximum number of server connections in the pool.
* `ssl_mode` - TLS mode of the connection pool.
* `updated_at` - Last update date of the connection pool.
* `uri` - Connection URI of the connection pool.
* `user_id` - ID of the user the pool connects with. Null when the pool uses a system user.
