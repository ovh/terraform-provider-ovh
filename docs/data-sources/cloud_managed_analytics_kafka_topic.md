---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_kafka_topic (Data Source)

Use this data source to get information about a topic of a kafka cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_managed_analytics_kafka_topic" "topic" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  id            = "ZZZ"
}

output "topic_name" {
  value = data.ovh_cloud_managed_analytics_kafka_topic.topic.name
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `id` - (Required) Topic ID

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `id` - See Argument Reference above.
* `min_insync_replicas` - Minimum insync replica accepted for this topic.
* `name` - Name of the topic.
* `partitions` - Number of partitions for this topic.
* `replication` - Number of replication for this topic.
* `retention_bytes` - Number of bytes for the retention of the data for this topic. Inferior to 0 mean Unlimited
* `retention_hours` - Number of hours for the retention of the data for this topic. Inferior to 0 mean Unlimited
* `service_name` - See Argument Reference above.
