data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_udp_farm" "farm_name" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name = "ingress-8080-gra"
  zone         = "gra"
  port         = 80
}

resource "ovh_iploadbalancing_udp_farm_server" "backend" {
  service_name           = "${data.ovh_iploadbalancing.lb.service_name}"
  farm_id                = "${ovh_iploadbalancing_udp_farm.farm_name.farm_id}"
  display_name           = "mybackend"
  address                = "4.5.6.7"
  status                 = "active"
  port                   = 80
}
