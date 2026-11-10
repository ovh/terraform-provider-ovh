data "ovh_cloud_managed_kubernetes_nodepool" "pool" {
  service_name          = "<Public cloud project id>"
  managed_kubernetes_id = "<cluster id>"
  id                    = "<nodepool id>"
}

output "ready_nodes" {
  value = data.ovh_cloud_managed_kubernetes_nodepool.pool.current_state.status.ready_nodes
}
