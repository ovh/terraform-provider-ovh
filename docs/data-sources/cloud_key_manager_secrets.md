---
subcategory : "Key Manager"
---

# ovh_cloud_key_manager_secrets (Data Source)

Use this data source to list all secrets in the Barbican Key Manager service for a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_key_manager_secrets" "all" {
  service_name = "Public cloud project ID"
}

output "secret_ids" {
  value = [for s in data.ovh_cloud_key_manager_secrets.all.secrets : s.id]
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

## Attributes Reference

The following attributes are exported:

* `secrets` - List of secrets. Each secret has the following attributes:
  * `id` - ID of the secret.
  * `checksum` - Computed hash representing the current resource state.
  * `created_at` - Creation date of the secret.
  * `updated_at` - Last update date of the secret.
  * `resource_status` - Secret readiness status.
  * `location` - Location of the secret:
    * `region` - Region of the secret.
    * `availability_zone` - Availability zone of the secret.
  * `name` - Name of the secret.
  * `secret_type` - Type of the secret.
  * `current_state` - Current state of the secret as reported by OpenStack Barbican.
