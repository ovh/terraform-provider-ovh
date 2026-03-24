---
subcategory : "Key Manager"
---

# ovh_cloud_keymanager_secret

Creates a secret in the Barbican Key Manager service for a public cloud project.

## Example Usage

Create an opaque secret with a base64-encoded payload:

```terraform
resource "ovh_cloud_keymanager_secret" "secret" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region       = "GRA"
  name         = "my-secret"
  secret_type  = "OPAQUE"

  payload              = base64encode("my-secret-value")
  payload_content_type = "APPLICATION_OCTET_STREAM"

  metadata = {
    environment = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `region` - (Required, Forces new resource) Region where the secret will be created.
* `name` - (Required, Forces new resource) Name of the secret.
* `secret_type` - (Required, Forces new resource) Type of the secret. Possible values: `SYMMETRIC`, `PUBLIC`, `PRIVATE`, `PASSPHRASE`, `CERTIFICATE`, `OPAQUE`.
* `algorithm` - (Optional, Forces new resource) Algorithm associated with the secret (e.g., `AES`, `RSA`).
* `bit_length` - (Optional, Forces new resource) Bit length of the secret (e.g., `256`).
* `mode` - (Optional, Forces new resource) Mode of the secret algorithm (e.g., `CBC`).
* `payload` - (Optional, Sensitive, Forces new resource) Secret payload data (base64-encoded). Write-only, never returned in responses. Requires `payload_content_type`.
* `payload_content_type` - (Optional, Forces new resource) Content type of the secret payload. Possible values: `TEXT_PLAIN`, `APPLICATION_OCTET_STREAM`, `APPLICATION_PKIX_CERT`, `APPLICATION_PKCS8`.
* `expiration` - (Optional, Forces new resource) Expiration date of the secret in RFC3339 format.
* `metadata` - (Optional) Key-value metadata for the secret. This is the only mutable field on a secret.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the secret.
* `checksum` - Computed hash representing the current resource state.
* `created_at` - Creation date of the secret.
* `updated_at` - Last update date of the secret.
* `resource_status` - Secret readiness status (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the secret as reported by OpenStack Barbican:
  * `name` - Name of the secret.
  * `secret_type` - Type of the secret.
  * `algorithm` - Algorithm of the secret.
  * `bit_length` - Bit length of the secret.
  * `mode` - Mode of the secret algorithm.
  * `payload_content_type` - Content type of the payload.
  * `expiration` - Expiration date of the secret.
  * `secret_ref` - OpenStack reference URL for the secret.
  * `status` - Status of the secret (`ACTIVE`, `ERROR`).
  * `region` - Region of the secret.
  * `metadata` - Key-value metadata.

## Import

A Key Manager secret can be imported using the `service_name` and `id`, separated by `/`:

```bash
$ terraform import ovh_cloud_keymanager_secret.my_secret service_name/secret_id
```
