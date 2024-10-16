---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_output_opensearch_index

Creates a DBaaS Logs Opensearch output index.

## Example Usage

```hcl
resource "ovh_dbaas_logs_output_opensearch_index" "index" {
  service_name = "...."
  description  = "my opensearch index"
  suffix = "index"
}
```

## Argument Reference

The following arguments are supported:
* `service_name` - (Required) The service name
* `description` - (Required) Index description
* `nb_shard` - (Required) Number of shards
* `suffix` - (Required) Index suffix

## Attributes Reference

Id is set to the opensearch index Id. In addition, the following attributes are exported:

* `alert_notify_enabled` - If set, notify when size is near 80, 90 or 100 % of its maximum capacity
* `created_at` - Index creation
* `current_size` - Current index size (in bytes)
* `description` - Index description
* `index_id` - Index ID
* `is_editable` - Indicates if you are allowed to edit entry
* `max_size` - Maximum index size (in bytes)
* `name` - Index name
* `nb_shard` - Number of shard
* `updated_at` - Index last update