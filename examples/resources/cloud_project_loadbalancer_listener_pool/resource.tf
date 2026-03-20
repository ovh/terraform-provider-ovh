resource "ovh_cloud_loadbalancer_listener_pool" "pool" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  listener_id     = ovh_cloud_loadbalancer_listener.http.id
  name            = "my-pool"
  protocol        = "HTTP"
  algorithm       = "ROUND_ROBIN"
}

resource "ovh_cloud_loadbalancer_listener_pool" "sticky" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  listener_id     = ovh_cloud_loadbalancer_listener.http.id
  name            = "sticky-pool"
  protocol        = "HTTP"
  algorithm       = "ROUND_ROBIN"

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "JSESSIONID"
  }
}
