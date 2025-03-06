---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_prometheus (Data Source)

Use this data source to get information about a prometheus of a database cluster associated with a public cloud project. For mongodb, please use ovh_cloud_project_database_mongodb_prometheus datasource

## Example Usage

```terraform
data "ovh_cloud_project_database_prometheus" "prometheus" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
}

output "name" {
  value = data.ovh_cloud_project_database_prometheus.prometheus.username
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required) The engine of the database cluster you want user information. To get a full list of available engine visit : [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases). Available engines:
  * `cassandra`
  * `kafka`
  * `kafkaConnect`
  * `kafkaMirrorMaker`
  * `mysql`
  * `opensearch`
  * `postgresql`
  * `redis`

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `engine` - See Argument Reference above.
* `id` - Cluster ID.
* `service_name` - See Argument Reference above.
* `targets` - List of all endpoint targets.
  * `Host` - Host of the endpoint.
  * `Port` - Connection port for the endpoint.
* `username` - name of the prometheus user.
