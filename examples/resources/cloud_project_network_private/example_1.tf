resource "ovh_cloud_project_network_private" "net" {
  service_name = "XXXXXX"
  name         = "admin_network"
  regions      = ["GRA11", "BHS5"]
}
