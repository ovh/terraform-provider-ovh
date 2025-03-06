resource "ovh_ip_move" "move_ip_to_load_balancer_xxxxx" {
  ip = "1.2.3.4"
  routed_to {
    service_name = "loadbalancer-XXXXX"
  }
}
