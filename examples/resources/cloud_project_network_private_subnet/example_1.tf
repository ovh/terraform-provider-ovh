resource "ovh_cloud_project_network_private_subnet" "subnet" {
  service_name = "xxxxx"
  network_id   = "0234543"
  region       = "GRA1"
  start        = "192.168.168.100"
  end          = "192.168.168.200"
  network      = "192.168.168.0/24"
  dhcp         = true
  no_gateway   = false
}
