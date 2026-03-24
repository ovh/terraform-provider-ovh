---
subcategory : "Key Manager"
---

# ovh_cloud_keymanager_secret_consumers (Data Source)

Use this data source to list all consumers registered on a Barbican Key Manager secret.

## Example Usage

```terraform
data "ovh_cloud_keymanager_secret_consumers" "consumers" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  secret_id    = "00000000-0000-0000-0000-000000000000"
}

output "consumer_services" {
  value = [for c in data.ovh_cloud_keymanager_secret_consumers.consumers.consumers : c.service]
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `secret_id` - (Required) UUID of the secret.

## Attributes Reference

The following attributes are exported:

* `consumers` - List of consumers registered on the secret. Each consumer has the following attributes:
  * `id` - Computed consumer identifier.
  * `service` - OpenStack service type of the consumer.
  * `resource_type` - Type of the resource consuming the secret.
  * `resource_id` - UUID of the resource consuming the secret.
