---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_role_permission

Reference a DBaaS logs role stream permission.

## Example Usage

```hcl
resource "ovh_dbaas_logs_role_permission" "permission" {
  service_name     = "ldp-xx-xxxxx"

  role_id = ovh_dbaas_logs_role.ro.id
  stream_id = ovh_dbaas_logs_output_graylog_stream.mystream.stream_id
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service name
* `role_id` -  (Required) The DBaaS Logs role id
* `stream_id` - (Required) The DBaaS Logs Graylog output stream id



