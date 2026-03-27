---
subcategory : "Key Manager"
---

# ovh_cloud_keymanager_containers (Data Source)

Use this data source to list all containers in the Barbican Key Manager service for a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_keymanager_containers" "all" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}

output "container_ids" {
  value = [for c in data.ovh_cloud_keymanager_containers.all.containers : c.id]
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

## Attributes Reference

The following attributes are exported:

* `containers` - List of containers. Each container has the following attributes:
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
    * `type` - Type of the container. Possible values: `CERTIFICATE`, `GENERIC`, `RSA`.
    * `container_ref` - OpenStack reference URL for the container.
    * `status` - Status of the container. Possible values: `ACTIVE`, `ERROR`.
    * `region` - Region of the container.
    * `secret_refs` - List of secret references:
      * `name` - Name of the secret reference.
      * `secret_id` - ID of the referenced secret.
