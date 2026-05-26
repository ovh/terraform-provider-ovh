data "ovh_vps_current_image" "current" {
  service_name = "vpsXXXXX.ovh.net"
}

output "current_image_name" {
  value = data.ovh_vps_current_image.current.name
}
