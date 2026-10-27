---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database

Creates a OVHcloud Managed Database Service in a public cloud project.

## Important

To learn more about OVHcloud Public Cloud Database please visit our [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).

The plans, flavors, number of nodes and regions available for each engine are listed in [Capabilities and Limitations of Public Cloud Databases](https://docs.ovhcloud.com/en/guides/public-cloud/databases/capabilities).

The capabilities documentation are helpful to find the properties of the configuration that matches your needs, such as plan, flavor, versions..

| Database |   Capabilities  |
| -------- | --------------- |
| MySQL | https://docs.ovhcloud.com/en/guides/public-cloud/databases/mysql-capabilities|
| PostgreSQL | https://docs.ovhcloud.com/en/guides/public-cloud/databases/postgresql-capabilities|
| MongoDB| https://docs.ovhcloud.com/en/guides/public-cloud/databases/mongodb-concept-capabilities|
| Valkey| https://docs.ovhcloud.com/en/guides/public-cloud/databases/valkey-capabilities|
| ClickHouse | https://docs.ovhcloud.com/en/guides/public-cloud/data-analytics/analytics/clickhouse-capabilities-limitations|
| Kafka | https://docs.ovhcloud.com/en/guides/public-cloud/data-analytics/analytics/capabilities|
| OpenSearch | https://docs.ovhcloud.com/en/guides/public-cloud/data-analytics/analytics/capabilities|


You can visit our dedicated Discord channel: https://discord.gg/PwPqWUpN8G. Ask questions, provide feedback and interact directly with the team that builds our databases services and terraform provider.

## Example Usage

Minimum settings for each engine (region choice is up to the user):

```terraform
resource "ovh_cloud_project_database" "kafkadb" {
  service_name          = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description           = "my-first-kafka"
  engine                = "kafka"
  version               = "4.3"
  flavor                = "b3-8"
  plan                  = "production"
  kafka_rest_api        = true
  kafka_schema_registry = true
  nodes {
    region = "DE"
  }
  nodes {
    region = "DE"
  }
  nodes {
    region = "DE"
  }
}

resource "ovh_cloud_project_database" "mongodb" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-first-mongodb"
  engine       = "mongodb"
  version      = "8.0"
  plan         = "production"
  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
  flavor = "db2-2"
}

resource "ovh_cloud_project_database" "mysqldb" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-first-mysql"
  engine       = "mysql"
  version      = "8.4"
  plan         = "discovery"
  nodes {
    region = "SBG"
  }
  flavor = "b3-4"
  advanced_configuration = {
    "mysql.sql_mode" : "ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE,STRICT_ALL_TABLES",
    "mysql.sql_require_primary_key" : "true"
  }
}

resource "ovh_cloud_project_database" "opensearchdb" {
  service_name            = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description             = "my-first-opensearch"
  engine                  = "opensearch"
  version                 = "3.3"
  plan                    = "discovery"
  opensearch_acls_enabled = true
  nodes {
    region = "UK"
  }
  flavor = "b3-4"
}

resource "ovh_cloud_project_database" "pgsqldb" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-first-postgresql"
  engine       = "postgresql"
  version      = "18"
  plan         = "discovery"
  nodes {
    region = "WAW"
  }
  flavor = "b3-4"
  ip_restrictions {
    description = "ip 1"
    ip          = "178.97.6.0/24"
  }
  ip_restrictions {
    description = "ip 2"
    ip          = "178.97.7.0/24"
  }
}

resource "ovh_cloud_project_database" "grafana" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-first-grafana"
  engine       = "grafana"
  version      = "11.6"
  plan         = "production"
  nodes {
    region = "GRA"
  }
  flavor = "b3-4"
}

resource "ovh_cloud_project_database" "valkeydb" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-first-valkey"
  engine       = "valkey"
  version      = "9.1"
  plan         = "discovery"
  nodes {
    region = "BHS"
  }
  flavor = "b3-4"
}
```

To deploy a production PostgreSQL service with two nodes on public network:

```terraform
resource "ovh_cloud_project_database" "postgresql" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-first-postgresql"
  engine       = "postgresql"
  version      = "18"
  plan         = "production"
  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
  flavor = "b3-8"
}
```

To deploy a production MongoDB service with three nodes on private network:

```terraform
resource "ovh_cloud_project_database" "mongodb" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-first-mongodb"
  engine       = "mongodb"
  version      = "8.0"
  plan         = "production"
  nodes {
    region     = "SBG"
    subnet_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  nodes {
    region     = "SBG"
    subnet_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  nodes {
    region     = "SBG"
    subnet_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  flavor = "db2-15"
}
```

### Number of nodes

The number of nodes of the cluster is the number of `nodes` blocks. Each plan requires a given number of nodes: for example, the MongoDB `production` plan requires exactly 3 nodes, and the `advanced` plan accepts 3 to 7 nodes. A wrong number of nodes fails at apply with a `NotEnoughNodes` error.

To deploy a production MongoDB service with three nodes:

```terraform
# The number of nodes is the number of "nodes" blocks: the production plan of
# MongoDB requires exactly 3 nodes. All nodes must be in the same region.
# Plans, flavors and number of nodes of MongoDB:
# https://docs.ovhcloud.com/en/guides/public-cloud/databases/mongodb-concept-capabilities
resource "ovh_cloud_project_database" "mongodb_production" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-production-mongodb"
  engine       = "mongodb"
  version      = "8.0"
  plan         = "production"
  flavor       = "db2-2"

  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
}
```

To deploy an advanced MongoDB service with five nodes:

```terraform
# The advanced plan of MongoDB accepts 3 to 7 nodes: this cluster has 5.
# All nodes must be in the same region.
# Plans, flavors and number of nodes of MongoDB:
# https://docs.ovhcloud.com/en/guides/public-cloud/databases/mongodb-concept-capabilities
resource "ovh_cloud_project_database" "mongodb_advanced" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-advanced-mongodb"
  engine       = "mongodb"
  version      = "8.0"
  plan         = "advanced"
  flavor       = "db2-15"

  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
}
```

### Network Update

You can update the network of an existing database service without recreating it.

To switch from public to private network, add `network_id` and `subnet_id` to the nodes:

```terraform
resource "ovh_cloud_project_database" "db" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-database"
  engine        = "postgresql"
  version       = "18"
  plan          = "production"
  nodes {
    region      = "GRA"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    subnet_id   = "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"
  }
  nodes {
    region      = "GRA"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    subnet_id   = "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"
  }
  flavor        = "b3-8"
}
```

To switch from private to public network, remove `network_id` and `subnet_id` from the nodes:

```terraform
resource "ovh_cloud_project_database" "db" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-database"
  engine        = "postgresql"
  version       = "18"
  plan          = "production"
  nodes {
    region      = "GRA"
  }
  nodes {
    region      = "GRA"
  }
  flavor        = "b3-8"
}
```

~> **Important:** Changing the network triggers a service rebuild. The service will go through `UPDATING` / `REBUILDING` states before returning to `RUNNING`. During this time the service may be temporarily unavailable. IP restrictions are cleared during a network change as the old IPs are no longer valid on the new network.

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `description` - (Optional) Small description of the database service.
* `engine` - (Required, Forces new resource) The database engine you want to deploy. To get a full list of available engine visit. [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).
* `flavor` - (Required) A valid OVHcloud public cloud database flavor name in which the nodes will be started. Ex: "b3-8". Changing this value upgrade the nodes with the new flavor. You can find the list of flavor names: https://www.ovhcloud.com/en/public-cloud/prices/
* `ip_restrictions` - (Optional) IP Blocks authorized to access to the cluster.
  * `description` - (Optional) Description of the IP restriction
  * `ip` - (Optional) Authorized IP
* `deletion_protection` - (Optional) Enable deletion protection
* `kafka_rest_api` - (Optional) Defines whether the REST API is enabled on a kafka cluster
* `kafka_schema_registry` - (Optional) Defines whether the schema registry is enabled on a Kafka cluster
* `nodes` - (Required, Minimum Items: 1) List of nodes object. The number of `nodes` blocks is the number of nodes of the cluster, and must match the number of nodes required by the plan. Multi region clusters are not available: all nodes must be identical.
  * `network_id` - (Optional) Private network id in which the node should be deployed. It's the regional openstackId of the private network. Can be updated in-place to change the network without recreating the service.
  * `region` - (Required) Public cloud region in which the node should be deployed. Ex: "GRA". Nodes can't be moved: changing the region of the nodes recreates the cluster.
  * `subnet_id` - (Optional) Private subnet ID in which the node is. Can be updated in-place to change the network without recreating the service.
* `opensearch_acls_enabled` - (Optional) Defines whether the ACLs are enabled on an OpenSearch cluster
* `disk_size` - (Optional) The disk size (in GB) of the database service.
* `advanced_configuration` - (Optional) Advanced configuration key / value.
* `plan` - (Required) Plan of the cluster. The plans available, and the number of nodes they require, depend on the engine (see [Capabilities and Limitations](https://docs.ovhcloud.com/en/guides/public-cloud/databases/capabilities)). For example:
  * ClickHouse: "discovery" (1 node), "production" (3 nodes).
  * MongoDB: "production" (3 nodes), "advanced" (3 to 7 nodes).
  * MySQL, PostgreSQL: "discovery" (1 node), "production" (2 nodes), "advanced" (3 nodes).
  * Valkey: "discovery" (1 node), "production" (2 nodes).
* `version` - (Required) The version of the engine in which the service should be deployed
* `backup_regions` - List of region where backups are pushed. Not more than 1 regions for MongoDB. Not more than 2 regions for the other engines with one being the same as the nodes[].region field
* `backup_time` - Time on which backups start every day (this parameter is not usable on the following engines: "grafana", "kafka", "kafkaconnect", "kafkamirrormaker", "opensearch").
* `maintenance_time` - Time on which maintenances can start every day. Format is: `hh:mm:ss` in UTC timezone, using 24 hours and zero padded values (e.g. `15:04:05`).

## Attributes Reference

The following attributes are exported:

* `id` - Public Cloud Database Service ID
* `service_name` - See Argument Reference above.
* `backup_regions` - See Argument Reference above.
* `backup_time` - See Argument Reference above.
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
* `ip_restrictions` - See Argument Reference above.
  * `description` - See Argument Reference above.
  * `ip` - See Argument Reference above.
  * `status` - Current status of the IP restriction.
* `deletion_protection` - Enable deletion protection
* `kafka_rest_api` - See Argument Reference above.
* `kafka_schema_registry` - See Argument Reference above.
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
* `disk_type` - Defines the disk type of the database service.
* `advanced_configuration` - Advanced configuration key / value.

## Timeouts

```terraform
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
