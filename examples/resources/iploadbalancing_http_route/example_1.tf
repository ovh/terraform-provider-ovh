resource "ovh_iploadbalancing_http_route" "https_redirect" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  display_name = "Redirect to HTTPS"
  weight = 1

  action {
    status = 302
    target = "https://$${host}$${path}$${arguments}"
    type   = "redirect"
  }
}
