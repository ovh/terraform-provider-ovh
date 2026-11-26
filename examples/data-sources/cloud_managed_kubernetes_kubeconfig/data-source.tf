data "ovh_cloud_managed_kubernetes_kubeconfig" "cluster" {
  service_name          = "<Public cloud project id>"
  managed_kubernetes_id = "<cluster id>"
}

resource "local_sensitive_file" "kubeconfig" {
  content  = data.ovh_cloud_managed_kubernetes_kubeconfig.cluster.content
  filename = "${path.module}/kubeconfig.yaml"
}
