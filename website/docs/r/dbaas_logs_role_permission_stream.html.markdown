---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_role_permission_stream

Reference a DBaaS logs role stream permission.

## Example Usage

```hcl
resource "ovh_dbaas_logs_role_permission_stream" "permission" {
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

## Attributes Reference
Id is set to the permission Id. In addition, the following attributes are exported:
* `permission_type` - Permission type (e.g., READ_ONLY)

## Import

DBaaS logs role stream permission can be imported using the `service_name`, `role_id` and `id`  of the permission, separated by "/" E.g.,

```bash
$  terraform import ovh_dbaas_logs_role_permission_stream.perm ldp-ra-XX/dc145bc2-eb01-4efe-a802-XXXXXX/e4818fa8-f426-11ef-a1f6-XXXXXXX
```

