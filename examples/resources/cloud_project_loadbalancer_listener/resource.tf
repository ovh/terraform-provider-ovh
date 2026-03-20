resource "ovh_cloud_loadbalancer_listener" "http" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  name            = "http-listener"
  protocol        = "HTTP"
  protocol_port   = 80
}

resource "ovh_cloud_loadbalancer_listener" "https" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  name            = "https-listener"
  protocol        = "TERMINATED_HTTPS"
  protocol_port   = 443
  description     = "HTTPS listener with TLS termination"

  default_tls_container_ref = "https://key-manager.cloud.ovh.net/v1/containers/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  tls_versions              = ["TLSv1.2", "TLSv1.3"]

  insert_headers {
    x_forwarded_for   = true
    x_forwarded_port  = true
    x_forwarded_proto = true
  }

  timeout_client_data    = 50000
  timeout_member_data    = 50000
  timeout_member_connect = 5000
}
