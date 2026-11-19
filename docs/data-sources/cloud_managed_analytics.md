---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics (Data Source)

Use this data source to get a managed analytics service (clickhouse, grafana, kafka, kafkaConnect, kafkaMirrorMaker, opensearch) of a public cloud project, using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_analytics" "service" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "kafka"
  id           = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
}
```

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `engine` - (Required) Engine of the service: clickhouse, grafana, kafka, kafkaConnect, kafkaMirrorMaker, opensearch.
* `id` - (Required) ID of the service.

## Attributes Reference

The following attributes are exported:

* `advanced_configuration` - Advanced configuration properties.
* `backup` - Backup configuration.
  * `regions` - Regions where the backups are stored.
  * `time` - Time of the daily backup, as HH:MM UTC.
* `checksum` - Hash of the service, computed by the API.
* `created_at` - Creation date of the service.
* `current_state` - Configuration running on the service.
  * `advanced_configuration` - Advanced configuration properties.
  * `backup` - Backup configuration.
    * `pitr` - Date until which point-in-time recovery is available.
    * `regions` - Regions where the backups are stored.
    * `retention_days` - Number of days the backups are kept.
    * `time` - Time of the daily backup, as HH:MM UTC.
  * `capabilities` - Capabilities of the service, such as USERS or CONNECTION_POOLS, with the state of each action: ENABLED or DISABLED.
    * `create` - State of the create action.
    * `delete` - State of the delete action.
    * `read` - State of the read action.
    * `update` - State of the update action.
  * `deletion_protected` - Whether the service is protected from deletion.
  * `description` - Description of the service.
  * `endpoints` - Connection endpoints of the service.
    * `component` - Component reached through the endpoint.
    * `domain` - Domain name of the endpoint.
    * `path` - Path of the endpoint.
    * `port` - Port of the endpoint.
    * `public` - Whether the endpoint uses the public domain name of the service.
    * `scheme` - Scheme of the endpoint.
    * `ssl` - Whether TLS is enabled on the endpoint.
    * `ssl_mode` - TLS mode of the endpoint.
    * `uri` - Connection URI of the endpoint.
  * `engine_capabilities` - Capabilities of the engine of the service.
    * `advanced_configuration_properties` - Advanced configuration properties that can be set, with their constraints.
      * `default` - Default value of the property.
      * `description` - Description of the property.
      * `maximum` - Maximum value of a numeric property.
      * `minimum` - Minimum value of a numeric property.
      * `name` - Name of the property.
      * `requires_restart` - Whether changing the property restarts the service.
      * `type` - Type of the property: BOOLEAN, DOUBLE, LONG or STRING.
      * `values` - Possible values of a string property.
    * `backup_regions` - Regions where the backups can be stored.
    * `categories` - Command categories a user can be given. Valkey only.
    * `commands` - Commands a user can be given. Valkey only.
    * `integrations` - Integrations the service can take part in.
      * `destination_engine` - Engine of the destination service.
      * `parameters` - Parameters of the integration.
        * `name` - Name of the parameter.
        * `type` - Type of the parameter: INTEGER or STRING.
      * `source_engine` - Engine of the source service.
      * `type` - Type of the integration.
  * `flavor` - Flavor of the nodes.
  * `maintenance_window` - Maintenance window.
    * `start_time_utc` - Start time of the maintenance window, as HH:MM UTC.
  * `network` - Network configuration.
    * `authorized_ips` - IP addresses allowed to connect to the service.
      * `rules` - Allowed CIDR blocks.
        * `cidr` - CIDR block.
        * `created_at` - Creation date of the rule.
        * `description` - Description of the rule.
    * `network_id` - ID of the private network the service is attached to.
    * `public` - Whether the service is published on a public address.
    * `subnet_id` - ID of the subnet the service is attached to.
  * `nodes` - Nodes of the service.
    * `created_at` - Creation date of the node.
    * `flavor` - Flavor of the node.
    * `id` - ID of the node.
    * `location` - Location of the node.
      * `availability_zone` - Availability zone of the node.
      * `region` - Region of the node.
    * `name` - Name of the node.
    * `port` - Port of the node.
    * `role` - Role of the node: PRIMARY, SECONDARY, STANDBY….
    * `status` - Readiness of the node.
  * `nodes_count` - Number of nodes.
  * `plan` - Plan of the service.
  * `prometheus_enabled` - Whether the Prometheus endpoint is enabled.
  * `status` - Status of the service reported by the engine.
  * `storage` - Storage configuration.
    * `size` - Storage size.
      * `unit` - Unit of the size: GB or MB.
      * `value` - Size, in the unit.
    * `type` - Type of storage.
  * `version` - Engine version.
* `current_tasks` - Asynchronous operations running on the service.
  * `errors` - Errors of the task.
    * `message` - Description of the error.
  * `id` - ID of the task.
  * `link` - Link to the details of the task.
  * `status` - Status of the task: ERROR, PENDING, RUNNING, SCHEDULED or WAITING_USER_INPUT.
  * `type` - Type of the task.
* `deletion_protected` - Whether the service is protected from deletion.
* `description` - Description of the service.
* `flavor` - Flavor of the nodes.
* `maintenance_window` - Maintenance window.
  * `start_time_utc` - Start time of the maintenance window, as HH:MM UTC.
* `network` - Network configuration.
  * `authorized_ips` - IP addresses allowed to connect to the service.
    * `rules` - Allowed CIDR blocks.
      * `cidr` - CIDR block.
      * `description` - Description of the rule.
  * `network_id` - ID of the private network the service is attached to.
  * `public` - Whether the service is published on a public address.
  * `subnet_id` - ID of the subnet the service is attached to.
* `nodes_count` - Number of nodes.
* `plan` - Plan of the service.
* `prometheus_enabled` - Whether the Prometheus endpoint is enabled.
* `region` - Region of the nodes.
* `resource_status` - Readiness of the service: CREATING, UPDATING, READY, ERROR, OUT_OF_SYNC, DELETING….
* `storage` - Storage configuration.
  * `size` - Storage size.
    * `unit` - Unit of the size: GB or MB.
    * `value` - Size, in the unit.
* `updated_at` - Last update date of the service.
* `version` - Engine version.
