---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_kafka_topics (Data Source)

Use this data source to get the list of topics of a kafka cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_kafka_topics" "topics" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "topic_ids" {
  value = data.ovh_cloud_project_database_kafka_topics.topics.topic_ids
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

`id` is set to the md5 sum of the list of all topics ids. In addition,
the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `topic_ids` - The list of topics ids of the kafka cluster associated with the project.
