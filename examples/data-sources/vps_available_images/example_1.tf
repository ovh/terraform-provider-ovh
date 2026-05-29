data "ovh_vps_available_images" "images" {
  service_name = "vpsXXXXX.ovh.net"
}

output "image_ids" {
  value = data.ovh_vps_available_images.images.images[*].id
}
