resource "ovh_cloud_project_region_network" "net" {
   service_name = "XXXXXX"
   region_name  = "EU-SOUTH-LZ-MAD-A"
   name         = "Madrid Network"
   subnet       = {
      cidr              = "10.0.0.0/24"
      enable_dhcp       = true
      enable_gateway_ip = false
      ip_version        = 4
   }
}
