resource "ovh_cloud_managed_kubernetes_nodepool" "pool" {
  service_name          = "<Public cloud project id>"
  managed_kubernetes_id = ovh_cloud_managed_kubernetes.cluster.id
  name                  = "my-pool"
  desired_nodes         = 3

  node_config = { flavor = "b3-8" }

  autoscaler = {
    enabled   = true
    min_nodes = 1
    max_nodes = 5
  }

  template = {
    metadata = { labels = { team = "backend" } }
  }

  # The autoscaler rewrites desired_nodes: keep the configured value as the initial size only.
  lifecycle {
    ignore_changes = [desired_nodes]
  }
}
