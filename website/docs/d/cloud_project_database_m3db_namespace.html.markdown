---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_m3db_namespace (Data Source)

Use this data source to get information about a namespace of a M3DB cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_m3db_namespace" "m3dbnamespace" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "m3dbnamespace_type" {
  value = data.ovh_cloud_project_database_m3db_namespace.m3dbnamespace.type
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `name` - (Required, Forces new resource) Name of the namespace.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `id` - ID of the namespace.
* `name` - See Argument Reference above.
* `resolution` - Resolution for an aggregated namespace.
* `retention_block_data_expiration_duration` - Controls how long we wait before expiring stale data.
* `retention_block_size_duration` - Controls how long to keep a block in memory before flushing to a fileset on disk.
* `retention_buffer_future_duration` - Controls how far into the future writes to the namespace will be accepted.
* `retention_buffer_past_duration` - Controls how far into the past writes to the namespace will be accepted.
* `retention_period_duration` - Controls the duration of time that M3DB will retain data for the namespace.
* `service_name` - See Argument Reference above.
* `snapshot_enabled`- SDefines whether M3db will create snapshot files for this namespace.
* `type` - Type of namespace.
* `writes_to_commit_log_enabled` - Defines whether M3DB will include writes to this namespace in the commit log.