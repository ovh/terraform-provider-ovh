data "ovh_cloud_managed_kubernetes_nodepool_nodes" "nodes" {
  service_name          = "<Public cloud project id>"
  managed_kubernetes_id = "<cluster id>"
  nodepool_id           = "<nodepool id>"
}

output "instance_ids" {
  value = [for n in data.ovh_cloud_managed_kubernetes_nodepool_nodes.nodes.nodes : n.current_state.instance_id if n.current_state.instance_id != null]
}
