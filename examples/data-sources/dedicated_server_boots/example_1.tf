data "ovh_dedicated_server_boots" "netboots" {
  service_name = "myserver"
  boot_type    = "harddisk"
}
