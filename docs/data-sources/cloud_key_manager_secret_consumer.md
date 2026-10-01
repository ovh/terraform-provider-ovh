---
subcategory : "Key Manager"
---

# ovh_cloud_key_manager_secret_consumer (Data Source)

Use this data source to get information about a single consumer registered on a Barbican Key Manager secret.

## Example Usage

```terraform
data "ovh_cloud_key_manager_secret_consumer" "consumer" {
  service_name = "Public cloud project ID"
  secret_id    = "00000000-0000-0000-0000-000000000000"
  consumer_id  = "Q09NUFVURTpJTlNUQU5DRToxMTExMTExMS0xMTExLTExMTEtMTExMS0xMTExMTExMTExMTE"
}

output "consumer_service" {
  value = data.ovh_cloud_key_manager_secret_consumer.consumer.service
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `secret_id` - (Required) UUID of the secret.
* `consumer_id` - (Required) Consumer identifier, as returned in the `id` attribute of the `ovh_cloud_key_manager_secret_consumer` resource and the `id` field of the `ovh_cloud_key_manager_secret_consumers` data source (the base64 identifier expected by the get-one endpoint).

## Attributes Reference

The following attributes are exported:

* `id` - Computed consumer identifier.
* `service` - OpenStack service type of the consumer.
* `resource_type` - Type of the resource consuming the secret.
* `resource_id` - UUID of the resource consuming the secret.
