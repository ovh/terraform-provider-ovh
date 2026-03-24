---
subcategory : "Key Manager"
---

# ovh_cloud_keymanager_container

Creates a container in the Barbican Key Manager service for a public cloud project. Containers hold references to secrets and are typically used to group related secrets such as TLS certificates.

## Example Usage

Create a certificate container referencing a certificate and a private key secret:

```terraform
resource "ovh_cloud_keymanager_secret" "certificate" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region       = "GRA"
  name         = "my-certificate"
  secret_type  = "CERTIFICATE"
}

resource "ovh_cloud_keymanager_secret" "private_key" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region       = "GRA"
  name         = "my-private-key"
  secret_type  = "PRIVATE"
}

resource "ovh_cloud_keymanager_container" "container" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region       = "GRA"
  name         = "my-certificate-container"
  type         = "CERTIFICATE"

  secret_refs {
    name      = "certificate"
    secret_id = ovh_cloud_keymanager_secret.certificate.id
  }

  secret_refs {
    name      = "private_key"
    secret_id = ovh_cloud_keymanager_secret.private_key.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `region` - (Required, Forces new resource) Region where the container will be created.
* `name` - (Required, Forces new resource) Name of the container.
* `type` - (Required, Forces new resource) Type of the container. Possible values: `CERTIFICATE`, `GENERIC`, `RSA`.
* `secret_refs` - (Optional) List of secret references in the container. Each `secret_refs` block supports:
  * `name` - (Required) Name of the secret reference (e.g., `certificate`, `private_key`, `public_key`).
  * `secret_id` - (Required) ID of the referenced secret.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the container.
* `checksum` - Computed hash representing the current resource state.
* `created_at` - Creation date of the container.
* `updated_at` - Last update date of the container.
* `resource_status` - Container readiness status (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the container as reported by OpenStack Barbican:
  * `name` - Name of the container.
  * `type` - Type of the container.
  * `container_ref` - OpenStack reference URL for the container.
  * `status` - Status of the container.
  * `region` - Region of the container.
  * `secret_refs` - List of secret references:
    * `name` - Name of the secret reference.
    * `secret_id` - ID of the referenced secret.

## Import

A Key Manager container can be imported using the `service_name` and `id`, separated by `/`:

```bash
$ terraform import ovh_cloud_keymanager_container.my_container service_name/container_id
```
