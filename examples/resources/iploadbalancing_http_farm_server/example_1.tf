data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_http_farm" "farmname" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  port         = 8080
  zone         = "all"
}

resource "ovh_iploadbalancing_http_farm_server" "backend" {
  service_name           = "${data.ovh_iploadbalancing.lb.service_name}"
  farm_id                = "${ovh_iploadbalancing_http_farm.farmname.id}"
  display_name           = "mybackend"
  address                = "4.5.6.7"
  status                 = "active"
  port                   = 80
  proxy_protocol_version = "v2"
  weight                 = 2
  probe                  = true
  ssl                    = false
  backup                 = true
}
