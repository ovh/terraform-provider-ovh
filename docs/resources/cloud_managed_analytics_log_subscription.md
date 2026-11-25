---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_log_subscription

Creates the forwarding of one kind of logs to a Logs Data Platform stream in a managed analytics service (clickhouse, grafana, kafka, kafkaConnect, kafkaMirrorMaker, opensearch), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "ldp-xx-xxxxx"
  title        = "my-stream"
}

resource "ovh_cloud_managed_analytics_log_subscription" "logs" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "kafka"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  stream_id    = data.ovh_dbaas_logs_output_graylog_stream.stream.id
  kind         = "customer_logs"
}
```

The API can't change a subscription: changing any argument recreates it.

The creation and the deletion wait for their operation on the Logs Data Platform service of the stream, which the provider reads with `GET /dbaas/logs/{serviceName}/operation/{operationId}`: the credentials need that right besides those on the project.

A creation adopts the subscription of the same kind to the same stream when there is one, such as one left by an interrupted apply, instead of forwarding the logs twice. When the wait for the operation fails, the subscription is still saved if it exists, and the next apply recreates it.

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used. Changing it recreates the log subscription.
* `engine` - (Required) Engine of the service: clickhouse, grafana, kafka, kafkaConnect, kafkaMirrorMaker, opensearch. Changing it recreates the log subscription.
* `service_id` - (Required) ID of the service. Changing it recreates the log subscription.
* `kind` - (Required) Kind of logs to forward. Available kinds are returned by the log/kind endpoint of the service. Changing it recreates the subscription.
* `stream_id` - (Required) ID of the Logs Data Platform stream receiving the logs. Changing it recreates the subscription.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the log subscription.
* `created_at` - Creation date of the subscription.
* `ldp_service_name` - Name of the Logs Data Platform service of the stream.
* `resource` - Resource the logs come from.
  * `name` - Name of the resource.
  * `type` - Type of the resource.
* `updated_at` - Last update date of the subscription.

## Timeouts

The creation and the deletion wait up to 20 minutes for the Logs Data Platform operation to end.

## Import

It can be imported using its `service_name`, `engine`, `service_id` and `id`, separated by "/".

```terraform
import {
  to = ovh_cloud_managed_analytics_log_subscription.logs
  id = "<service_name>/<engine>/<service_id>/<id>"
}
```

```bash
$ terraform import ovh_cloud_managed_analytics_log_subscription.logs service_name/engine/service_id/id
```
