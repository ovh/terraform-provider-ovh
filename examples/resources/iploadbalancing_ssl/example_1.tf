data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_ssl" "sslname" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name = "test"
  certificate  = "..."
  key          = "..."
  chain        = "..."

  # use this if ssl is configured as frontend default_ssl
  lifecycle {
    create_before_destroy = true
  }
}
