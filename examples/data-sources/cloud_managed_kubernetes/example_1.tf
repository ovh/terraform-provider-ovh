data "ovh_cloud_managed_kubernetes" "my_kube_cluster" {
  service_name = "XXXXXX"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "version" {
  value = data.ovh_cloud_managed_kubernetes.my_kube_cluster.version
}

output "kubeconfig" {
  value = data.ovh_cloud_managed_kubernetes.my_kube_cluster.kubeconfig
  sensitive = true
}

output "kube_host" {
  value = data.ovh_cloud_managed_kubernetes.my_kube_cluster.kubeconfig_attributes[0].host
}