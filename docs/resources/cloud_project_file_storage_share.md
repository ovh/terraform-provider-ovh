---
subcategory: "Cloud Project"
---

# ovh_cloud_project_file_storage_share (Resource)

Creates a file storage share in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_project_file_storage_share" "share" {
  service_name = "xxxxxxxxxx"
  region_name  = "GRA11"
  name         = "my_share"
  description  = "My file storage share"
  size         = 150
  type         = "standard-1az"
  network_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subnet_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The ID of the public cloud project.
* `region_name` - (Required, Forces new resource) The region in which the share will be created.
* `description` - (Optional) Share description.
* `name` - (Optional) Share name.
* `size` - (Optional) Share size in Gigabytes.
* `type` - (Optional, Forces new resource) Share type. Currently only `standard-1az` is supported.
* `network_id` - (Required, Forces new resource) Private network ID.
* `subnet_id` - (Required, Forces new resource) Subnet ID.
* `snapshot_id` - (Optional, Forces new resource) Snapshot ID used to create the share.
* `availability_zone` - (Optional, Forces new resource) Availability zone of the share (required in 3AZ regions).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the share (UUID).
* `created_at` - Share creation date.
* `is_public` - Whether the share is public.
* `protocol` - Share protocol (e.g. `NFS`).
* `share_network_id` - Share network ID.
* `status` - Share status.

## Import

A file storage share can be imported using the `service_name`, `region_name`, and `share_id` separated by `/`:

```bash
$ terraform import ovh_cloud_project_file_storage_share.share service_name/region_name/share_id
```
