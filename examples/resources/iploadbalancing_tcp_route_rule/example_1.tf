resource "ovh_iploadbalancing_tcp_route" "reject" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  weight       = 1
  frontend_id  = 11111

  action {
    type = "reject"
  }
}

resource "ovh_iploadbalancing_tcp_route_rule" "example_rule" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  route_id     = ovh_iploadbalancing_tcp_route.reject.id
  display_name = "Match example.com host"
  field        = "sni"
  match        = "is"
  negate       = false
  pattern      = "example.com"
}
