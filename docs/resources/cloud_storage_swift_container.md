---
subcategory: "Object Storage"
---

# ovh_cloud_storage_swift_container

Creates a Swift object storage container in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_storage_swift_container" "container" {
  service_name = <Public cloud project id>
  name         = "my-container"
  region       = "GRA"
}
```

A container hosting a static website, reachable from a browser on `https://www.example.com`:

```terraform
resource "ovh_cloud_storage_swift_container" "website" {
  service_name = <Public cloud project id>
  name         = "my-website"
  region       = "GRA"
  type         = "STATIC"
  cors         = ["https://www.example.com"]
}
```

A container using the Public Cloud Archive storage policy. The policy can only be chosen when the container is created:

```terraform
resource "ovh_cloud_storage_swift_container" "archive" {
  service_name = <Public cloud project id>
  name         = "my-archive"
  region       = "GRA"
  archive      = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `name` - (Required) Container name, unique within the region: 1 to 256 characters, must not contain a slash (`/`). **Changing this value recreates the resource.**
* `region` - (Required) Region identifier hosting the container. Must be upper-case (e.g. `GRA`, `SBG`, `BHS`). **Changing this value recreates the resource.**
* `archive` - (Optional) Create the container with the Public Cloud Archive storage policy instead of the standard policy. Objects in an archive container are sealed after upload and must be unsealed before download. Defaults to `false`. **Changing this value recreates the resource.**
* `type` - (Optional) Access type of the container. One of:
  * `PRIVATE` (default) - Objects need a token or a presigned URL to be read.
  * `PUBLIC` - Anyone can read and list the objects.
  * `STATIC` - Public, plus static website hosting: `index.html` is served as the index page, `error.html` on errors and `listing.css` styles the listings. Default versions of these three files are pushed to the container when it is switched to `STATIC`.

  Switching back to `PRIVATE` removes the public read ACL and the website settings but keeps the objects.
* `cors` - (Optional) Set of origins allowed to access the container from a browser (`Access-Control-Allow-Origin` container metadata), e.g. `https://www.example.com`. Defaults to no origin, which blocks cross-origin browser requests.

## Attributes Reference

The following attributes are exported:

* `id` - Container identifier, in the `<REGION>_<name>` form (e.g. `GRA_my-container`).
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
  * `objects_count` - Number of objects stored in the container.
  * `objects_size` - Total size of the stored objects, in bytes.
  * `static_url` - Public URL of the container. Only set when `type` is `PUBLIC` or `STATIC`.

## Import

A Swift container can be imported using the `service_name` and the container `id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_storage_swift_container.container
  id = "<service_name>/<container_id>"
}
```

~> __NOTE__ The container `id` is `<REGION>_<name>` (for example `GRA_my-container`).

```bash
$ terraform import ovh_cloud_storage_swift_container.container service_name/container_id
```
