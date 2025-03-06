resource "ovh_iploadbalancing_http_route_rule" "example_rule" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  route_id     = "${ovh_iploadbalancing_http_route.https_redirect.id}"
  display_name = "Match example.com Host header"
  field        = "headers"
  match        = "is"
  negate       = false
  pattern      = "example.com"
  sub_field    = "Host"
}
