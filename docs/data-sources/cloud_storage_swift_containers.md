---
subcategory: "Object Storage"
---

# ovh_cloud_storage_swift_containers (Data Source)

List the Swift object storage containers in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_storage_swift_containers" "containers" {
  service_name = <Public cloud project id>
}
```

Filter the containers by region:

```terraform
data "ovh_cloud_storage_swift_containers" "containers" {
  service_name = <Public cloud project id>
  region       = "GRA"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `region` - (Optional) If set, only containers located in this region are returned.

## Attributes Reference

* `containers` - List of containers:
  * `id` - Container identifier, in the `<REGION>_<name>` form.
  * `name` - Container name.
  * `location` - Target region of the container:
    * `region` - Region identifier.
  * `archive` - Whether the container uses the Public Cloud Archive storage policy.
  * `type` - Access type of the container (`PRIVATE`, `PUBLIC`, `STATIC`).
  * `cors` - Set of origins allowed to access the container from a browser.
  * `checksum` - Computed hash representing the current target specification value.
  * `created_at` - Creation date of the container.
  * `updated_at` - Last update date of the container.
  * `resource_status` - Container readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `SUSPENDED`, `UNKNOWN`, `UPDATING`).
  * `current_state` - Current observed state of the container, read from the Swift backend:
    * `name` - Container name.
    * `location` - Region hosting the container:
      * `region` - Region identifier.
    * `archive` - Whether the container uses the Public Cloud Archive storage policy.
    * `type` - Access type (`PRIVATE`, `PUBLIC`, `STATIC`).
    * `cors` - Origins currently allowed.
    * `objects_count` - Number of objects stored in the container. May be `null` on a list read.
    * `objects_size` - Total size of the stored objects, in bytes. May be `null` on a list read.
    * `static_url` - Public URL of the container. Only set when `type` is `PUBLIC` or `STATIC`, may be `null` on a list read.
