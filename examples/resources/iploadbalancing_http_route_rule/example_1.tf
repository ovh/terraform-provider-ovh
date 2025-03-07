resource "ovh_iploadbalancing_http_route" "https_redirect" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  display_name = "Redirect to HTTPS"
  weight       = 1
  frontend_id  = 11111

  action {
    status = 302
    target = "https://$${host}$${path}$${arguments}"
    type   = "redirect"
  }
}

resource "ovh_iploadbalancing_http_route_rule" "example_rule" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  route_id     = "${ovh_iploadbalancing_http_route.https_redirect.id}"
  display_name = "Match example.com host"
  field        = "host"
  match        = "is"
  negate       = false
  pattern      = "example.com"
}
