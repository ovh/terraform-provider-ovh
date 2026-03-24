---
subcategory : "Key Manager"
---

# ovh_cloud_keymanager_container (Data Source)

Use this data source to get information about a single container in the Barbican Key Manager service.

## Example Usage

```terraform
data "ovh_cloud_keymanager_container" "container" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  container_id = "00000000-0000-0000-0000-000000000000"
}

output "container_name" {
  value = data.ovh_cloud_keymanager_container.container.name
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `container_id` - (Required) ID of the container.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the container.
* `checksum` - Computed hash representing the current resource state.
* `created_at` - Creation date of the container.
* `updated_at` - Last update date of the container.
* `resource_status` - Container readiness status.
* `region` - Region of the container.
* `name` - Name of the container.
* `type` - Type of the container.
* `current_state` - Current state of the container as reported by OpenStack Barbican:
  * `name` - Name of the container.
  * `type` - Type of the container.
  * `container_ref` - OpenStack reference URL.
  * `status` - Status of the container.
  * `region` - Region of the container.
  * `secret_refs` - List of secret references:
    * `name` - Name of the secret reference.
    * `secret_id` - ID of the referenced secret.
