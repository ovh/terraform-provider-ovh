---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_ip_restrictions (Data Source)

Deprecated: Use ip_restrictions field in cloud_project_database datasource instead.

Use this data source to get the list of IP restrictions associated with a public cloud project.

## Example Usage

To get the list of IP restriction on a database cluster service:

```hcl
data "ovh_cloud_project_database_ip_restrictions" "ip_restrictions" {
  service_name  = "XXXXXX"
  engine        = "YYYY"
  cluster_id    = "ZZZZ"
}

output "ips" {
  value = data.ovh_cloud_project_database_ip_restrictions.ip_restrictions.ips
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required) The engine of the database cluster you want to list IP restrictions. To get a full list of available engine visit:
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).

* `cluster_id` - (Required) Cluster ID


## Attributes Reference

`id` is set to the md5 sum of the list of all IP restrictions. In addition,
the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `engine` - See Argument Reference above.
* `ips` - The list of IP restriction of the database associated with the project.
* `service_name` - See Argument Reference above.
