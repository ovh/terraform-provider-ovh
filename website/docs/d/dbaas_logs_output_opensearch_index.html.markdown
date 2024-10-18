---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_output_opensearch_index (Data Source)

Use this data source to retrieve information about a DBaas logs output opensearch index.

## Example Usage

```hcl

data "ovh_dbaas_logs_output_opensearch_index" "index" {
  service_name = "ldp-xx-xxxxx"
  name        = "index-name"
}
```

## Argument Reference

* `service_name` - The service name. It's the ID of your Logs Data Platform instance.
* `name` - Index name

## Attributes Reference

`id` is set to output opensearch index ID. In addition, the following attributes are exported:

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