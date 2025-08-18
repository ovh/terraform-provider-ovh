---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_output_opensearch_alias

Creates a DBaaS Logs Opensearch output alias.

## Example Usage

```terraform
resource "ovh_dbaas_logs_output_opensearch_alias" "alias" {
  service_name = "...."
  description  = "my opensearch alias"
  suffix = "alias"
}
```

## Argument Reference

The following arguments are supported:
* `service_name` - (Required) The service name
* `description` - (Required) Index description
* `suffix` - (Required) Index suffix
* `indexes` - (Optional) List of attached indexes id
* `streams` - (Optional) List of attached streams id

## Attributes Reference

Id is set to the opensearch alias Id. In addition, the following attributes are exported:

* `alias_id` - Alias Id
* `created_at` - Alias creation
* `description` - Alias description
* `indexes` - List of attached indexes id
* `is_editable` - Indicates if you are allowed to edit entry
* `max_size` - Maximum index size (in bytes)
* `name` - Alias name
* `nb_index` - Number of indices linked
* `nb_stream` - Number of streams linked
* `streams` - List of attached streams id
* `updated_at` - Input last update

## Import

DBaaS logs Opensearch output alias can be imported using the `service_name` and `alias_id`, separated by "/" E.g.,

```bash
$ terraform import ovh_dbaas_logs_output_opensearch_alias.this ldp-XX-NNNNN/7e2301b6-c0fe-41d9-bc58-44fd844f1e3e
```
