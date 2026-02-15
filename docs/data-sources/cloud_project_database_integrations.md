---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_integrations (Data Source)

Use this data source to get the list of integrations of a database cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_database_integrations" "integrations" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
}

output "integration_ids" {
  value = data.ovh_cloud_project_database_integrations.integrations.integration_ids
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required) The engine of the database cluster you want to list integrations. To get a full list of available engine visit: [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases). All engines available except `mongodb`

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

`id` is set to the md5 sum of the list of all integration ids. In addition, the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `engine` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `integration_ids` - The list of integrations ids of the database cluster associated with the project.
