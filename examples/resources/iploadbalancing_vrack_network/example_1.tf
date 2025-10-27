data "ovh_iploadbalancing" "iplb" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
}

resource "ovh_vrack_iploadbalancing" "vip_lb" {
  service_name     = "xxx"
  ip_loadbalancing = data.ovh_iploadbalancing.iplb.service_name
}

resource "ovh_iploadbalancing_vrack_network" "network" {
  service_name = ovh_vrack_iploadbalancing.vip_lb.ip_loadbalancing
  subnet       = "10.0.0.0/16"
  vlan         = 1
  nat_ip       = "10.0.0.0/27"
  display_name = "mynetwork"
}

resource "ovh_iploadbalancing_tcp_farm" "test_farm" {
  service_name     = ovh_iploadbalancing_vrack_network.network.service_name
  display_name     = "mytcpbackends"
  port             = 80
  vrack_network_id = ovh_iploadbalancing_vrack_network.network.vrack_network_id
  zone             = tolist(data.ovh_iploadbalancing.iplb.zone)[0]
}
