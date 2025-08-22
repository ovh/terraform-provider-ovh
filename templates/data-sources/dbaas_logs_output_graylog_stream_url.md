---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_output_graylog_stream_url (Data Source)

Use this data source to retrieve the list of URLs for a DBaas logs output Graylog stream.

## Example Usage

```terraform
data "ovh_dbaas_logs_output_graylog_stream_url" "urls" {
  service_name = "ldp-xx-xxxxx"
  stream_id    = "STREAM_ID"
}
```

## Argument Reference

* `service_name` - The service name. It's the ID of your Logs Data Platform instance.
* `stream_id` - Stream ID.

## Attributes Reference

The following attributes are exported:

* `url` - List of URLs. Each element contains:
  * `address` - URL address
  * `type` - URL type (e.g. `GRAYLOG_WEBUI`, `WEB_SOCKET`)