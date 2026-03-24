---
subcategory : "Key Manager"
---

# ovh_cloud_keymanager_secret (Data Source)

Use this data source to get information about a single secret in the Barbican Key Manager service.

## Example Usage

```terraform
data "ovh_cloud_keymanager_secret" "secret" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  secret_id    = "00000000-0000-0000-0000-000000000000"
}

output "secret_name" {
  value = data.ovh_cloud_keymanager_secret.secret.name
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `secret_id` - (Required) ID of the secret.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the secret.
* `checksum` - Computed hash representing the current resource state.
* `created_at` - Creation date of the secret.
* `updated_at` - Last update date of the secret.
* `resource_status` - Secret readiness status.
* `region` - Region of the secret.
* `name` - Name of the secret.
* `secret_type` - Type of the secret.
* `metadata` - Key-value metadata for the secret.
* `current_state` - Current state of the secret as reported by OpenStack Barbican:
  * `name` - Name of the secret.
  * `secret_type` - Type of the secret.
  * `algorithm` - Algorithm of the secret.
  * `bit_length` - Bit length of the secret.
  * `mode` - Mode of the secret algorithm.
  * `payload_content_type` - Content type of the payload.
  * `expiration` - Expiration date.
  * `secret_ref` - OpenStack reference URL.
  * `status` - Status of the secret.
  * `region` - Region of the secret.
  * `metadata` - Key-value metadata.
