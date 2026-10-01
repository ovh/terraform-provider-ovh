---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_output_graylog_stream

Creates a DBaaS Logs Graylog output stream.

## Example Usage

### Example 1 - Basic stream

```terraform
resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "...."
  title        = "my stream"
  description  = "my graylog stream"
}
```

### Example 2 - Stream retention

```terraform
data "ovh_dbaas_logs_cluster_retention" "retention" {
  service_name = "ldp-xx-xxxxx"
  cluster_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  duration     = "P14D"
}

resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "...."
  title        = "my stream"
  description  = "my graylog stream"
  retention_id = data.ovh_dbaas_logs_cluster_retention.retention.retention_id
}
```

### Example 3 - Cold Storage encryption using a PGP public key

```terraform
resource "ovh_dbaas_logs_encryption_key" "key" {
  service_name = "ldp-xx-xxxxx"
  title        = "my PGP key"
  content      = file("my-pgp-public-key.asc")
  fingerprint  = "ABCDEF1234567890ABCDEF1234567890ABCDEF12"
}

resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name            = "ldp-xx-xxxxx"
  title                   = "my stream"
  description             = "my encrypted graylog stream"
  cold_storage_enabled    = true
  cold_storage_target     = "PCA"
  cold_storage_retention  = 1
  encryption_keys_ids     = [ovh_dbaas_logs_encryption_key.key.id]
}
```

## Argument Reference

The following arguments are supported:
* `service_name` - (Required) The service name
* `title` - (Required) Stream name
* `description` - (Required) Stream description
* `parent_stream_id` - Parent stream ID
* `retention_id` - Retention ID
* `cold_storage_compression` - Cold storage compression method. One of "LZMA", "GZIP", "DEFLATED", "ZSTD"
* `cold_storage_content` - ColdStorage content. One of "ALL", "GELF", "PLAIN"
* `cold_storage_enabled` - Is Cold storage enabled?
* `cold_storage_notify_enabled` - Notify on new Cold storage archive
* `cold_storage_retention` - Cold storage retention in year
* `cold_storage_target` - ColdStorage destination. One of "PCA", "PCS"
* `encryption_keys_ids` - Set of encryption key IDs used to encrypt stream archives
* `indexing_enabled` - Enable ES indexing
* `indexing_max_size` - Maximum indexing size (in GB)
* `indexing_notify_enabled` - If set, notify when size is near 80, 90 or 100 % of the maximum configured setting
* `pause_indexing_on_max_size` - If set, pause indexing when maximum size is reached
* `web_socket_enabled` - Enable Websocket

## Attributes Reference

Id is set to the output stream Id. In addition, the following attributes are exported:

* `can_alert` - Indicates if the current user can create alert on the stream
* `created_at` - Stream creation
* `is_editable` - Indicates if you are allowed to edit entry
* `is_shareable` - Indicates if you are allowed to share entry
* `nb_alert_condition` - Number of alert condition
* `nb_archive` - Number of coldstored archives
* `stream_id` - Stream ID
* `updated_at` - Stream last update
* `write_token` - Write token of the stream (empty if the caller is not the owner of the stream)

## Import

DBaas logs output Graylog stream can be imported using the `service_name` of the cluster and `stream_id` of the graylog output stream, separated by "/" E.g.,

```bash
$ terraform import ovh_dbaas_logs_output_graylog_stream.ldp ldp-az-12345/9d2f9cf8-9f92-1337-c0f3-48a0213d2c6f
```