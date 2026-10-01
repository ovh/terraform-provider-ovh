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
