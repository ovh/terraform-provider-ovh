---
subcategory: "Cloud Project"
---

# ovh_cloud_project_file_storage_share_network (Resource)

Creates a share network in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_project_file_storage_share_network" "sn" {
  service_name = "xxxxxxxxxx"
  region_name  = "GRA"
  name         = "my-share-network"
  description  = "Shared network for file storage"
  network_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subnet_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  # availability_zone = "eu-west-par-a" # required in 3AZ regions
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The ID of the public cloud project.
* `region_name` - (Required, Forces new resource) The region in which the share network will be created.
* `network_id` - (Required, Forces new resource) Private network ID.
* `subnet_id` - (Required, Forces new resource) Subnet ID.
* `name` - (Optional, Forces new resource) Share network name.
* `description` - (Optional, Forces new resource) Share network description.
* `availability_zone` - (Optional, Forces new resource) Availability zone of the share network (required in 3AZ regions).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the share network (UUID).
* `cidr` - Subnet CIDR inherited from the Neutron subnet.
* `network_type` - Share network type.
* `created_at` - Share network creation date.
* `updated_at` - Share network last update date.

## Import

A share network can be imported using the `service_name`, `region_name`, and `share_network_id` separated by `/`:

```bash
$ terraform import ovh_cloud_project_file_storage_share_network.sn service_name/region_name/share_network_id
```
