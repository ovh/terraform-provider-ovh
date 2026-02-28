---
subcategory : "Managed Databases"
---

~> **DEPRECATED:** Use `ovh_cloud_managed_database_m3db_namespace` instead. This resource will be removed in the next major version.

# ovh_cloud_project_database_m3db_namespace

Creates a namespace for a M3DB cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_database" "m3db" {
  service_name  = "XXX"
  engine        = "m3db"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_m3db_namespace" "namespace" {
  service_name              = data.ovh_cloud_project_database.m3db.service_name
  cluster_id                = data.ovh_cloud_project_database.m3db.id
  name                      = "mynamespace"
  resolution                = "P2D"
  retention_period_duration = "PT48H"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `name` - (Required, Forces new resource) Name of the namespace. A namespace named "default" is mapped with already created default namespace instead of creating a new namespace.

* `resolution` - (Optional) Resolution for an aggregated namespace. Should follow Rfc3339 e.g P2D, PT48H.

* `retention_block_data_expiration_duration` - (Optional) Controls how long we wait before expiring stale data. Should follow Rfc3339 e.g P2D, PT48H.

* `retention_block_size_duration` - (Optional, Forces new resource) Controls how long to keep a block in memory before flushing to a fileset on disk. Should follow Rfc3339 e.g P2D, PT48H.

* `retention_buffer_future_duration` - (Optional) Controls how far into the future writes to the namespace will be accepted. Should follow Rfc3339 e.g P2D, PT48H.

* `retention_buffer_past_duration` - (Optional) Controls how far into the past writes to the namespace will be accepted. Should follow Rfc3339 e.g P2D, PT48H.

* `retention_period_duration` - (Optional) Controls the duration of time that M3DB will retain data for the namespace. Should follow Rfc3339 e.g P2D, PT48H.

* `snapshot_enabled` - (Optional) Defines whether M3DB will create snapshot files for this namespace.

* `writes_to_commit_log_enabled` - (Optional) Defines whether M3DB will include writes to this namespace in the commit log.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `id` - ID of the namespace.
* `name` - See Argument Reference above.
* `resolution` - See Argument Reference above.
* `retention_block_data_expiration_duration` - See Argument Reference above.
* `retention_block_size_duration` - See Argument Reference above.
* `retention_buffer_future_duration` - See Argument Reference above.
* `retention_buffer_past_duration` - See Argument Reference above.
* `retention_period_duration` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `snapshot_enabled`- See Argument Reference above.
* `type` - Type of namespace.
* `writes_to_commit_log_enabled` - See Argument Reference above.

## Timeouts

```terraform
resource "ovh_cloud_project_database_m3db_namespace" "namespace" {
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

OVHcloud Managed M3DB clusters namespaces can be imported using the `service_name`, `cluster_id` and `id` of the namespace, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_database_m3db_namespace.my_namespace service_name/cluster_id/id
```
