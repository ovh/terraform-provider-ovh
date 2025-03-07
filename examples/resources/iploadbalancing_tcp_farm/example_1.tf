data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_tcp_farm" "farm_name" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name = "ingress-8080-gra"
  zone         = "GRA"
}
