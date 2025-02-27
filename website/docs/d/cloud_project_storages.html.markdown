---
subcategory : "Object Storage"
---

# ovh_cloud_project_storages

List your S3™* compatible storage container.
\* S3 is a trademark filed by Amazon Technologies,Inc. OVHcloud's service is not sponsored by, endorsed by, or otherwise affiliated with Amazon Technologies,Inc.

## Example Usage

```hcl
data "ovh_cloud_project_storage" "storage" {
  service_name = "<public cloud project ID>"
  region_name = "GRA"
}
```

## Schema

### Required

- `region_name` (String) Region name
- `service_name` (String) Service name

### Read-Only

- `containers` (Attributes Set) (see [below for nested schema](#nestedatt--containers))

<a id="nestedatt--containers"></a>
### Nested Schema for `containers`

Read-Only:

- `created_at` (String) The date and timestamp when the resource was created
- `encryption` (Attributes) Encryption configuration (see [below for nested schema](#nestedatt--containers--encryption))
- `name` (String) Container name
- `objects` (Attributes List) Container objects (see [below for nested schema](#nestedatt--containers--objects))
- `objects_count` (Number) Container total objects count
- `objects_size` (Number) Container total objects size (bytes)
- `owner_id` (Number) Container owner user ID
- `region` (String) Container region
- `replication` (Attributes) Replication configuration (see [below for nested schema](#nestedatt--containers--replication))
- `tags` (Map of String) Container tags
- `versioning` (Attributes) Versioning configuration (see [below for nested schema](#nestedatt--containers--versioning))
- `virtual_host` (String) Container virtual host

<a id="nestedatt--containers--encryption"></a>
### Nested Schema for `containers.encryption`

Read-Only:

- `sse_algorithm` (String) Encryption algorithm


<a id="nestedatt--containers--objects"></a>
### Nested Schema for `containers.objects`

Read-Only:

- `etag` (String) ETag
- `is_delete_marker` (Boolean) Whether this object is a delete marker
- `is_latest` (Boolean) Whether this is the latest version of the object
- `key` (String) Key
- `last_modified` (String) Last modification date
- `size` (Number) Size (bytes)
- `storage_class` (String) Storage class
- `version_id` (String) Version ID of the object


<a id="nestedatt--containers--replication"></a>
### Nested Schema for `containers.replication`

Read-Only:

- `rules` (Attributes List) Replication rules (see [below for nested schema](#nestedatt--containers--replication--rules))

<a id="nestedatt--containers--replication--rules"></a>
### Nested Schema for `containers.replication.rules`

Read-Only:

- `delete_marker_replication` (String) Delete marker replication
- `destination` (Attributes) Rule destination configuration (see [below for nested schema](#nestedatt--containers--replication--rules--destination))
- `filter` (Attributes) Rule filters (see [below for nested schema](#nestedatt--containers--replication--rules--filter))
- `id` (String) Rule ID
- `priority` (Number) Rule priority
- `status` (String) Rule status

<a id="nestedatt--containers--replication--rules--destination"></a>
### Nested Schema for `containers.replication.rules.destination`

Read-Only:

- `name` (String) Destination bucket name
- `region` (String) Destination region, can be null if destination bucket has been deleted
- `storage_class` (String) Destination storage class


<a id="nestedatt--containers--replication--rules--filter"></a>
### Nested Schema for `containers.replication.rules.filter`

Read-Only:

- `prefix` (String) Prefix filter
- `tags` (Map of String) Tags filter

<a id="nestedatt--containers--versioning"></a>
### Nested Schema for `containers.versioning`

Read-Only:

- `status` (String) Versioning status
