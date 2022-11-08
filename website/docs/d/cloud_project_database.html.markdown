---
layout: "ovh"
page_title: "OVH: cloud_project_database"
sidebar_current: "docs-ovh-datasource-cloud-project-database-x"
description: |-
  Get information about a managed database cluster in a public cloud project.
---

# ovh_cloud_project_database (Data Source)

Use this data source to get the managed database of a public cloud project.

## Example Usage

To get information of a database cluster service:
```hcl
data "ovh_cloud_project_database" "db" {
  service_name  = "XXXXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

output "cluster_id" {
  value = data.ovh_cloud_project_database.db.id
}
```

## Argument Reference


* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required) The database engine you want to get information. To get a full list of available engine visit:
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).

* `id` - (Required) Cluster ID


## Attributes Reference


The following attributes are exported:

* `id` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `backup_time` - Time on which backups start every day.
* `created_at` - Date of the creation of the cluster.
* `description` - Small description of the database service.
* `endpoints` - List of all endpoints objects of the service.
  * `component` - Type of component the URI relates to.
  * `domain` - Domain of the cluster.
  * `path` - Path of the endpoint.
  * `port` - Connection port for the endpoint.
  * `scheme` - Scheme used to generate the URI.
  * `ssl` - Defines whether the endpoint uses SSL.
  * `ssl_mode` - SSL mode used to connect to the service if the SSL is enabled.
  * `uri` - URI of the endpoint.
* `engine` - See Argument Reference above.
* `flavor` - A valid OVHcloud public cloud database flavor name in which the nodes will be started.
* `kafka_rest_api` - Defines whether the REST API is enabled on a kafka cluster.
* `maintenance_time` - Time on which maintenances can start every day.
* `network_type` - Type of network of the cluster.
* `nodes` - List of nodes object.
  * `network_id` - Private network id in which the node should be deployed. It's the regional openstackId of the private network
  * `region` - Public cloud region in which the node should be deployed.
  * `subnet_id` -  Private subnet ID in which the node is.
* `plan` - Plan of the cluster.
* `status` - Current status of the cluster.
* `version` - The version of the engine in which the service should be deployed
* `disk_size` - The disk size (in GB) of the database service.
* `disk_type` -  The disk type of the database service.
