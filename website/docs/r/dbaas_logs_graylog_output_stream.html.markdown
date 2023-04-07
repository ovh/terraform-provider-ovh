---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_graylog_output_stream

Creates a dbaas logs graylog output stream.

## Example Usage

```hcl
resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "...."
  title        = "my stream"
  description  = "my graylog stream"
}
```

## Argument Reference

The following arguments are supported:
* `service_name` - (Required) The service name
* `description` - (Required) Stream description
* `title` - (Required) Stream description
* `parent_stream_id` - Parent stream ID
* `retention_id` - Retention ID
* `cold_storage_compression` - Cold storage compression method. One of "LZMA", "GZIP", "DEFLATED", "ZSTD"
* `cold_storage_content` - ColdStorage content. One of "ALL", "GLEF", "PLAIN"
* `cold_storage_enabled` - Is Cold storage enabled?
* `cold_storage_notify_enabled` - Notify on new Cold storage archive
* `cold_storage_retention` - Cold storage retention in year
* `cold_storage_target` - ColdStorage destination. One of "PCA", "PCS"
* `indexing_enabled` - Enable ES indexing
* `indexing_max_size` - Maximum indexing size (in GB)
* `indexing_notify_enabled` - If set, notify when size is near 80, 90 or 100 % of the maximum configured setting
* `pause_indexing_on_max_size` - If set, pause indexing when maximum size is reach
* `web_socket_enabled` - Enable Websocket

## Attributes Reference

Id is set to the output stream Id. In addition, the following attributes are exported:

* `can_alert` - Indicates if the current user can create alert on the stream
* `created_at` - Stream creation
* `is_editable` - Indicates if you are allowed to edit entry
* `is_shareable` - Indicates if you are allowed to share entry
* `nb_alert_condition` - Number of alert condition
* `nb_archive` - Number of coldstored archivesr
* `stream_id` - Stream ID
* `updated_at` - Stream last updater
