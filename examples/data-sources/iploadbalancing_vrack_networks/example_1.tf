data "ovh_iploadbalancing_vrack_networks" "lb_networks" {
  service_name = "XXXXXX"
  subnet       = "10.0.0.0/24"
}
