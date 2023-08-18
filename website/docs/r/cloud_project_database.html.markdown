---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database

Creates a OVHcloud Managed Database Service in a public cloud project.

## Important

To learn more about OVHcloud Public Cloud Database please visit our 
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).


You can visit our dedicated Discord channel: https://discord.gg/PwPqWUpN8G. Ask questions, provide feedback and 
interact directly with the team that builds our databases services and terraform provider.

## Example Usage

Minimum settings for each engine (region choice is up to the user):

```hcl
resource "ovh_cloud_project_database" "cassandradb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-cassandra"
  engine        = "cassandra"
  version       = "4.0"
  plan          = "essential"
  nodes {
    region  = "BHS"
  }
  nodes {
    region  = "BHS"
  }
  nodes {
    region  = "BHS"
  }
  flavor        = "db1-4"
}

resource "ovh_cloud_project_database" "kafkadb" {
  service_name    = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description     = "my-first-kafka"
  engine          = "kafka"
  version         = "3.4"
  plan            = "business"
  kafka_rest_api  = true
  nodes {
    region  = "DE"
  }
  nodes {
    region  = "DE"
  }
  nodes {
    region  = "DE"
  }
  flavor           = "db1-4"
}

resource "ovh_cloud_project_database" "m3db" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-m3db"
  engine        = "m3db"
  version       = "1.2"
  plan          = "essential"
  nodes {
    region  = "BHS"
  }
  flavor        = "db1-7"
}

resource "ovh_cloud_project_database" "mongodb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-mongodb"
  engine        = "mongodb"
  version       = "5.0"
  plan          = "essential"
  nodes {
    region =  "GRA"
  }
  flavor        = "db1-2"
}

resource "ovh_cloud_project_database" "mysqldb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-mysql"
  engine        = "mysql"
  version       = "8"
  plan          = "essential"
  nodes {
    region  = "SBG"
  }
  flavor        = "db1-4"
  advanced_configuration = {
    "mysql.sql_mode": "ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE,STRICT_ALL_TABLES",
    "mysql.sql_require_primary_key": "true"
  }
}

resource "ovh_cloud_project_database" "opensearchdb" {
  service_name            = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description             = "my-first-opensearch"
  engine                  = "opensearch"
  version                 = "1"
  plan                    = "essential"
  opensearch_acls_enabled = true
  nodes {
    region = "UK"
  }
  flavor                  = "db1-4"
}

resource "ovh_cloud_project_database" "pgsqldb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-postgresql"
  engine        = "postgresql"
  version       = "14"
  plan          = "essential"
  nodes {
    region  = "WAW"
  }
  flavor        = "db1-4"
}

resource "ovh_cloud_project_database" "redisdb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-redis"
  engine        = "redis"
  version       = "6.2"
  plan          = "essential"
  nodes {
    region  = "BHS"
  }
  flavor        = "db1-4"
}

resource "ovh_cloud_project_database" "grafana" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-grafana"
  engine        = "grafana"
  version       = "9.1"
  plan          = "essential"
  nodes {
    region =  "GRA"
  }
  flavor        = "db1-4"
}
```

To deploy a business PostgreSQL service with two nodes on public network:

```hcl
resource "ovh_cloud_project_database" "postgresql" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-postgresql"
  engine        = "postgresql"
  version       = "14"
  plan          = "business"
  nodes {
    region  = "GRA"
  }
  nodes {
    region  = "GRA"
  }
  flavor        = "db1-15"
}
```


To deploy an enterprise MongoDB service with three nodes on private network:

```hcl
resource "ovh_cloud_project_database" "mongodb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-mongodb"
  engine        = "mongodb"
  version       = "5.0"
  plan          = "enterprise"
  nodes {
    region      = "SBG"
    subnet_id   = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  nodes {
    region      = "SBG"
    subnet_id   = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  nodes {
    region      = "SBG"
    subnet_id   = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  flavor        = "db1-30"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `description` - (Optional) Small description of the database service.

* `engine` - (Required, Forces new resource) The database engine you want to deploy. To get a full list of available engine visit.
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).

* `flavor` -  (Required) A valid OVHcloud public cloud database flavor name in which the nodes will be started.
  Ex: "db1-7". Changing this value upgrade the nodes with the new flavor.
  You can find the list of flavor names: https://www.ovhcloud.com/fr/public-cloud/prices/

* `kafka_rest_api` -  (Optional) Defines whether the REST API is enabled on a kafka cluster

* `nodes` - (Required, Minimum Items: 1) List of nodes object.
  Multi region cluster are not yet available, all node should be identical.
  * `network_id` - (Optional, Forces new resource) Private network id in which the node should be deployed. It's the regional openstackId of the private network
  * `region` - (Required, Forces new resource) Public cloud region in which the node should be deployed.
    Ex: "GRA'.
  * `subnet_id` - (Optional, Forces new resource) Private subnet ID in which the node is.

* `opensearch_acls_enabled` -  (Optional) Defines whether the ACLs are enabled on an OpenSearch cluster

* `disk_size` -  (Optional) The disk size (in GB) of the database service.

* `advanced_configuration` -  (Optional) Advanced configuration key / value.

* `plan` - (Required) Plan of the cluster.
  Enum: "essential", "business", "enterprise".

* `version` - (Required) The version of the engine in which the service should be deployed

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
* `kafka_rest_api` - See Argument Reference above.
* `maintenance_time` - Time on which maintenances can start every day.
* `network_type` - Type of network of the cluster.
* `nodes` - See Argument Reference above.
  * `network_id` - See Argument Reference above.
  * `region` - See Argument Reference above.
  * `subnet_id` - See Argument Reference above.
* `opensearch_acls_enabled` - See Argument Reference above.
* `plan` - See Argument Reference above.
* `status` - Current status of the cluster.
* `version` - See Argument Reference above.
* `disk_size` - See Argument Reference above.
* `disk_type` -  Defines the disk type of the database service.
* `advanced_configuration` - Advanced configuration key / value.



## Timeouts

```hcl
resource "ovh_cloud_project_database" "db" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}
```
* `create` - (Default 20m)
* `update` - (Default 40m)
* `delete` - (Default 20m)

## Import

OVHcloud Managed database clusters can be imported using the `service_name`, `engine`, `id` of the cluster, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_database.my_database_cluster service_name/engine/id
```
