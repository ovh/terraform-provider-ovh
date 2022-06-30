---
layout: "ovh"
page_title: "OVH: cloud_project_database"
sidebar_current: "docs-ovh-resource-cloud-project-database-x"
description: |-
Creates a managed database in a public cloud project.
---

# ovh_cloud_project_database

Creates a OVH Managed Database Service in a public cloud project.

## Important

This resource is in beta state, you should use it with care.


To learn more about OVHcloud Public Cloud Database please visit our 
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).


You can visit our dedicated Discord channel: https://discord.gg/PwPqWUpN8G. Ask questions, provide feedback and 
interact directly with the team that builds our databases services and terraform provider.

## Example Usage

To deploy a business PostgreSQL service with two nodes on public network:
```hcl
resource "ovh_cloud_project_database" "postgresql" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-first-postgresql"
  engine       = "postgresql"
  version      = "14"
  plan         = "business"
  nodes {
    region     = "GRA"
  }
  nodes {
    region     = "GRA"
  }
  flavor = "db1-15"
}
```


To deploy an enterprise MongoDB service with three nodes on private network:
```hcl
resource "ovh_cloud_project_database" "postgresql" {
  service_name = var.openstack_infos.project_id
  description  = "my-first-postgresql"
  engine       = "mongodb"
  version      = "5.0"
  plan         = "enterprise"
  nodes {
    region     = "SBG"
    subnet_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"  
  }
  nodes {
    region     = "SBG"
    subnet_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  nodes {
    region     = "SBG"
    subnet_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  flavor = "db1-30"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `description` - (Optional) Small description of the database service.

* `engine` - The database engine you want to deploy. To get a full list of available engine visit.
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).

* `flavor` -  a valid OVH public cloud database flavor name in which the nodes will be started.
  Ex: "db1-7". Changing this value upgrade the nodes with the new flavor.
  You can find the list of flavor names: https://www.ovhcloud.com/fr/public-cloud/prices/

* `nodes` - List of nodes object.
  Multi region cluster are not yet available, all node should be identical.
  * `network_id` - Private network id in which the node should be deployed.
  * `region` - Public cloud region in which the node should be deployed.
    Ex: "GRA'.
  * `subnet_id` - Private subnet ID in which the node is.

* `plan` - List of nodes object.
  Enum: "esential", "business", "enterprise".

* `version` - The version of the engine in which the service should be deployed

## Attributes Reference

The following attributes are exported:

* `id` - Public Cloud Database Service ID
* `service_name` - See Argument Reference above.
* `backup_time` - Time on which backups start every day.
* `created_at` - Date of the creation of the cluster.
* `description` - See Argument Reference above.
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
* `flavor` - See Argument Reference above.
* `maintenance_time` - Time on which maintenances can start every day.
* `network_type` - Type of network of the cluster.
* `nodes` - See Argument Reference above.
* `plan` - See Argument Reference above.
* `status` - Current status of the cluster.
* `version` - See Argument Reference above.