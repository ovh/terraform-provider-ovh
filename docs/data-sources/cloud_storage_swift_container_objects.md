---
subcategory: "Object Storage"
---

# ovh_cloud_storage_swift_container_objects (Data Source)

List the objects stored in a Swift object storage container of a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_storage_swift_container_objects" "objects" {
  service_name = <Public cloud project id>
  container_id = "GRA_my-container"
}
```

Only the objects under the `documents/` prefix, 100 at most:

```terraform
data "ovh_cloud_storage_swift_container_objects" "documents" {
  service_name = <Public cloud project id>
  container_id = "GRA_my-container"
  prefix       = "documents/"
  limit        = 100
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `container_id` - (Required) Container identifier, in the `<REGION>_<name>` form (e.g. `GRA_my-container`).
* `prefix` - (Optional) If set, only objects whose key starts with this prefix are returned.
* `limit` - (Optional) Maximum number of objects returned. When unset, every object is returned, paging through the container.

## Attributes Reference

* `objects` - List of objects:
  * `key` - Object key, its full path within the container (e.g. `documents/report.pdf`).
  * `size` - Object size, in bytes.
  * `content_type` - MIME type of the object.
  * `last_modified_at` - Last modification date of the object.
  * `etag` - MD5 hash of the object content, or a hash of the segment hashes for a large object.
  * `retrieval_state` - Retrieval state of the object, only meaningful in an archive container (`SEALED`, `UNSEALING`, `UNSEALED`). `null` in a standard container.
  * `retrieval_delay` - Number of seconds before the object becomes readable. Only set while `retrieval_state` is `UNSEALING`.
