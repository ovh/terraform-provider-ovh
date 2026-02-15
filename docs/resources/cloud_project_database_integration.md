---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_integration

Creates an integration for a database cluster associated with a public cloud project.

With this resource you can create an integration for all engine except `mongodb`.

Please take a look at the list of available `types` in the `Argument references` section in order to know the list of available integrations. For example, thanks to the integration feature you can have your PostgreSQL logs in your OpenSearch Database.

## Example Usage

Push PostgreSQL logs in an OpenSearch DB:

```terraform
data "ovh_cloud_project_database" "db_postgresql" {
  service_name  = "XXXX"
  engine        = "postgresql"
  id            = "ZZZZ"
}

data "ovh_cloud_project_database" "db_opensearch" {
  service_name  = "XXXX"
  engine        = "opensearch"
  id            = "ZZZZ"
}

resource "ovh_cloud_project_database_integration" "integration" {
  service_name            = data.ovh_cloud_project_database.db_postgresql.service_name
  engine                  = data.ovh_cloud_project_database.db_postgresql.engine
  cluster_id              = data.ovh_cloud_project_database.db_postgresql.id
  source_service_id       = data.ovh_cloud_project_database.db_postgresql.id
  destination_service_id  = data.ovh_cloud_project_database.db_opensearch.id
  type                    = "opensearchLogs"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required, Forces new resource) The engine of the database cluster you want to add. You can find the complete list of available engine in the [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases). All engines available except `mongodb`.

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `destination_service_id` - (Required, Forces new resource) ID of the destination service.

* `parameters` - (Optional, Forces new resource) Parameters for the integration.

* `source_service_id` - (Required, Forces new resource) ID of the source service.

* `type` - (Optional, Forces new resource) Type of the integration. Available types:
  * `grafanaDashboard`
  * `grafanaDatasource`
  * `kafkaConnect`
  * `kafkaLogs`
  * `kafkaMirrorMaker`
  * `m3aggregator`
  * `m3dbMetrics`
  * `opensearchLogs`
  * `postgresqlMetrics`

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `destination_service_id` - See Argument Reference above.
* `engine` - See Argument Reference above.
* `id` - ID of the integration.
* `parameters` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `source_service_id` - See Argument Reference above.
* `status` - Current status of the integration.
* `type` - See Argument Reference above.

## Timeouts

```terraform
resource "ovh_cloud_project_database_integration" "integration" {
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

OVHcloud Managed database cluster integrations can be imported using the `service_name`, `engine`, `cluster_id` and `id` of the integration, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_database_integration.my_integration service_name/engine/cluster_id/id
```
