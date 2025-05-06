---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_log_subscription

Creates a log subscription for a cluster associated with a public cloud project.

## Example Usage

Create a log subscription for a database.

```terraform
data "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "ldp-xx-xxxxx"
  title        = "my stream"
}

data "ovh_cloud_project_database" "db" {
  service_name  = "XXX"
  engine        = "YYY"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_log_subscription" "subscription" {
	service_name = data.ovh_cloud_project_database.db.service_name
	engine       = data.ovh_cloud_project_database.db.engine
	cluster_id   = data.ovh_cloud_project_database.db.id
	stream_id    = data.ovh_dbaas_logs_output_graylog_stream.stream.id
  kind         = "customer_logs"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `engine` - (Required, Forces new resource) The database engine for which you want to manage a subscription. To get a full list of available engine visit. [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).
* `cluster_id` - (Required, Forces new resource) Cluster ID.
* `stream_id` - (Required, Forces new resource) Id of the target Log data platform stream.
* `kind` - (Required, Forces new resource) Log kind name of this subscription.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Creation date of the subscription.
* `engine` - See Argument Reference above.
* `id` - ID of the log subscription.
* `kind` - See Argument Reference above.
* `ldp_service_name` - Name of the destination log service.
* `operation_id` - Identifier of the operation.
* `resource_name` - Name of subscribed resource, where the logs come from.
* `resource_type` - Type of subscribed resource, where the logs come from.
* `service_name` - See Argument Reference above.
* `stream_id` - See Argument Reference above.
* `updated_at` - Last update date of the subscription.

## Timeouts

```terraform
resource "ovh_cloud_project_database_log_subscription" "sub" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}
```
* `create` - (Default 20m)
* `update` - (Default 20m)
* `delete` - (Default 20m)

## Import

OVHcloud Managed clusters logs subscription can be imported using the `service_name`, `engine`, `cluster_id` and `id` of the subscription, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_database_log_subscription.sub service_name/engine/cluster_id/id
```
