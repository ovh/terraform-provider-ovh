data "ovh_vps_available_image" "image" {
  service_name = "vpsXXXXX.ovh.net"
  id           = "debian12_64_virtio_202409"
}

output "image_name" {
  value = data.ovh_vps_available_image.image.name
}
