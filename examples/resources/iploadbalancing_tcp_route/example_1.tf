resource "ovh_iploadbalancing_tcp_route" "tcp_reject" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  weight = 1

  action {
    type = "reject"
  }
}
