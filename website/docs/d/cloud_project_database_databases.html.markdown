---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_databases (Data Source)

Use this data source to get the list of databases of a database cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_databases" "databases" {
  service_name  = "XXXX"
  engine        = "YYYY"
  cluster_id    = "ZZZ"
}

output "database_ids" {
  value = data.ovh_cloud_project_database_databases.databases.database_ids
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required) The engine of the database cluster you want to list databases. To get a full list of available engine visit:
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).
Available engines:
  * `mysql`
  * `postgresql`

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

`id` is set to the md5 sum of the list of all database ids. In addition,
the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `database_ids` - The list of databases ids of the database cluster associated with the project.
* `engine` - See Argument Reference above.
* `service_name` - See Argument Reference above.
