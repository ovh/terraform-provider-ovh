data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_tcp_farm" "farm80" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name = "ingress-8080-gra"
  zone         = "all"
  port         = 80
}

resource "ovh_iploadbalancing_tcp_frontend" "test_frontend" {
  service_name    = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name    = "ingress-8080-gra"
  zone            = "all"
  port            = "80,443"
  default_farm_id = "${ovh_iploadbalancing_tcp_farm.farm80.id}"
}
