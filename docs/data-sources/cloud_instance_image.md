---
subcategory: "Instances"
---

# ovh_cloud_instance_image (Data Source)

Use this data source to retrieve information about an image available in a public cloud project. This is read-only reference data: it describes a bootable image the project can create instances from, not an existing resource.

## Example Usage

```terraform
data "ovh_cloud_instance_image" "image" {
  service_name = "<Public cloud project id>"
  id           = "<image id>"
}

# Minimum sizing a flavor must provide to boot this image.
output "image_requirements" {
  value = format(
    "%s requires at least %d GB of disk and %d MB of RAM",
    data.ovh_cloud_instance_image.image.name,
    data.ovh_cloud_instance_image.image.min_disk,
    data.ovh_cloud_instance_image.image.min_ram,
  )
}
```

To look an image up by name instead of by id, use
[`ovh_cloud_instance_images`](cloud_instance_images.md).

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `id` - (Required) The OpenStack/Glance image ID to look up. Stable within a region and used to reference the image when creating an instance.

## Attributes Reference

The following attributes are exported:

* `name` - Display name of the image as reported by the backend (for example the distribution and version, such as `Debian 12`).
* `status` - Availability status of the image as reported by the backend. Only images in an active status can be used to create an instance.
* `visibility` - Visibility scope of the image, for example whether it is a public OVHcloud-provided image or private to the project.
* `min_disk` - Minimum root disk size, in GB, that an instance must provide to boot from this image. A flavor whose disk is smaller than this value cannot be used with the image.
* `min_ram` - Minimum amount of memory, in MB, that an instance must provide to boot from this image. A flavor whose RAM is below this value cannot be used with the image.
* `size` - Size of the image on the backend, expressed in bytes.
* `created_at` - Timestamp at which the image was created on the backend, in RFC 3339 format.
* `updated_at` - Timestamp of the last modification of the image on the backend, in RFC 3339 format.
* `location` - Region (and, where applicable, availability zone) where this image is offered. The image catalog is per-region: an image returned for one region is not guaranteed to exist in another:
  * `region` - Region code.
  * `availability_zone` - Availability zone within the region.
