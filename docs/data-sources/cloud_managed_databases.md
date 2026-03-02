---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_databases (Data Source)

Use this data source to get the list of managed databases of a public cloud project.

## Example Usage

To get the list of database clusters service for a given engine:

```terraform
data "ovh_cloud_managed_databases" "dbs" {
  service_name  = "XXXXXX"
  engine        = "YYYY"
}

output "cluster_ids" {
  value = data.ovh_cloud_managed_databases.dbs.cluster_ids
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required) The database engine you want to list. To get a full list of available engine visit: [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).

## Attributes Reference

`id` is set to the md5 sum of the list of all databases ids. In addition, the following attributes are exported:

* `cluster_ids` - The list of managed databases ids of the project.
* `engine` - See Argument Reference above.
* `service_name` - See Argument Reference above.
