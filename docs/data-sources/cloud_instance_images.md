---
subcategory: "Instances"
---

# ovh_cloud_instance_images (Data Source)

Use this data source to list the images available in a public cloud project. This is read-only reference data: each entry describes a bootable image the project can create instances from, not an existing resource.

## Example Usage

```terraform
# Omit `region` to list the catalog of every region at once: the same image
# then appears once per region it is offered in.
data "ovh_cloud_instance_images" "images" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
}

# Resolve an image id by name, to feed ovh_cloud_instance.image_id.
output "debian_12_image_id" {
  value = one([
    for image in data.ovh_cloud_instance_images.images.images :
    image.id if image.name == "Debian 12"
  ])
}

# Names of every image the project can currently boot from in the region.
output "bootable_image_names" {
  value = sort([
    for image in data.ovh_cloud_instance_images.images.images :
    image.name if lower(image.status) == "active"
  ])
}
```

### Feeding an instance

```terraform
data "ovh_cloud_instance_images" "catalog" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
}

resource "ovh_cloud_instance" "instance" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-instance"
  flavor_id    = "<flavor id>"

  image_id = one([
    for image in data.ovh_cloud_instance_images.catalog.images :
    image.id if image.name == "Debian 12"
  ])

  networks = [
    { auto_assign_public_ip = true },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `region` - (Optional) Restrict the listing to the images offered in this region. The catalog is per-region: an image returned for one region is not guaranteed to exist in another.

## Attributes Reference

The following attributes are exported:

* `images` - Images advertised by the backend for this project, one entry per region the image is offered in:
  * `id` - The OpenStack/Glance image ID. Stable within a region and used to reference the image when creating an instance.
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
