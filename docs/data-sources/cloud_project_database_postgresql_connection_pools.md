---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_postgresql_connection_pools (Data Source)

Use this data source to get the list of connection pools of a postgresql cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_database_postgresql_connection_pools" "test_pools" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "connection_pool_ids" {
  value = data.ovh_cloud_project_database_postgresql_connection_pools.test_pools.connection_pool_ids
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `cluster_id` - (Required) Cluster ID.

## Attributes Reference

`id` is set to the md5 sum of the list of all patterns ids. In addition, the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `connection_pool_ids` - The list of patterns ids of the opensearch cluster associated with the project.
