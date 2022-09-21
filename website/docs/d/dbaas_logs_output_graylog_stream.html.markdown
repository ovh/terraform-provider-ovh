---
layout: "ovh"
page_title: "OVH: dbaas_logs_output_graylog_stream"
sidebar_current: "docs-ovh-datasource-dbaas-logs-output-graylog-stream"
description: |-
  Get information & status of a DBaas logs output graylog stream.
---

# ovh_dbaas_logs_output_graylog_stream (Data Source)

Use this data source to retrieve information about a DBaas logs output graylog stream.

## Example Usage

```hcl

data "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "XXXXXX"
  title        = "my stream"
}
```

## Argument Reference

* `service_name` - The service name
* `title` - Stream description

## Attributes Reference

`id` is set to output graylog stream ID. In addition, the following attributes are exported:

* `cold_storage_compression` - Cold storage compression method
* `cold_storage_content` - ColdStorage content
* `cold_storage_enabled` - Is Cold storage enabled?
* `cold_storage_notify_enabled` - Notify on new Cold storage archive
* `cold_storage_retention` - Cold storage retention in year
* `cold_storage_target` - ColdStorage destination
* `created_at` - Stream creation
* `description` - Stream description
* `indexing_enabled` - Enable ES indexing
* `indexing_max_size` - Maximum indexing size (in GB)
* `indexing_notify_enabled` - If set, notify when size is near 80, 90 or 100 % of the maximum configured setting
* `is_editable` - Indicates if you are allowed to edit entry
* `is_shareable` - Indicates if you are allowed to share entry
* `nb_alert_condition` - Number of alert condition
* `nb_archive` - Number of coldstored archives
* `parent_stream_id` - Parent stream ID
* `pause_indexing_on_max_size` - If set, pause indexing when maximum size is reach
* `retention_id` - Retention ID
* `stream_id` - Stream ID
* `updated_at` - Stream last update
* `web_socket_enabled` - Enable Websocket
