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
