---
layout: "ovh"
page_title: "OVH: cloud_project_database_kafka_topic"
sidebar_current: "docs-ovh-resource-cloud-project-database-kafka-topic"
description: |-
  Creates a topic for a kafka cluster associated with a public cloud project.
---

# ovh_cloud_project_database_kafka_topic

Creates a topic for a kafka cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database" "kafka" {
  service_name  = "XXX"
  engine        = "kafka"
  cluster_id    = "ZZZ"
}

resource "ovh_cloud_project_database_kafka_topic" "topic" {
	service_name = ovh_cloud_project_database.kafka.service_name
	cluster_id   = ovh_cloud_project_database.kafka.id
	name = "mytopic"
	min_insync_replicas = 1
	partitions = 3
	replication = 2
	retention_bytes = 4
	retention_hours = 5
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `name` - (Required, Forces new resource) Name of the topic. No spaces allowed.

* `min_insync_replicas` - (Optional, Forces new resource) Minimum insync replica accepted for this topic. Should be superior to 0

* `partitions` - (Optional, Forces new resource) Number of partitions for this topic. Should be superior to 0

* `replication` - (Optional, Forces new resource) Number of replication for this topic. Should be superior to 1

* `retention_bytes` - (Optional, Forces new resource) Number of bytes for the retention of the data for this topic. Inferior to 0 means unlimited

* `retention_hours` - (Optional, Forces new resource) Number of hours for the retention of the data for this topic. Should be superior to -2. Inferior to 0 means unlimited



## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `id` - ID of the topic.
* `min_insync_replicas` - See Argument Reference above.
* `name` - See Argument Reference above.
* `partitions` - See Argument Reference above.
* `replication` - See Argument Reference above.
* `retention_bytes` - See Argument Reference above.
* `retention_hours` - See Argument Reference above.
* `service_name` - See Argument Reference above.

## Import

OVHcloud Managed kafka clusters topics can be imported using the `service_name`, `cluster_id` and `id` of the topic, separated by "/" E.g.,

```
$ terraform import ovh_cloud_project_database_kafka_topic.my_topic <service_name>/<cluster_id>/<id>