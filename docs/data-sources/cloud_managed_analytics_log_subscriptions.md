---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_log_subscriptions (Data Source)

Use this data source to list the IDs of the log subscriptions of a managed analytics service (clickhouse, grafana, kafka, kafkaConnect, kafkaMirrorMaker, opensearch), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_analytics_log_subscriptions" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "kafka"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  kind         = "customer_logs"
}
```

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `engine` - (Required) Engine of the service: clickhouse, grafana, kafka, kafkaConnect, kafkaMirrorMaker, opensearch.
* `service_id` - (Required) ID of the service.
* `kind` - (Optional) Lists only the subscriptions of this kind of logs.

## Attributes Reference

The following attributes are exported:

* `subscription_ids` - The IDs of the log subscriptions.
