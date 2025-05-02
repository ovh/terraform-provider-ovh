resource "ovh_cloud_project_kube" "my_cluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA11"
}

output "my_cluster_host" {
  value = ovh_cloud_project_kube.my_cluster.kubeconfig_attributes[0].host
  sensitive = true
}

output "my_cluster_cluster_ca_certificate" {
  value = ovh_cloud_project_kube.my_cluster.kubeconfig_attributes[0].cluster_ca_certificate
  sensitive = true
}

output "my_cluster_client_certificate" {
  value = ovh_cloud_project_kube.my_cluster.kubeconfig_attributes[0].client_certificate
  sensitive = true
}

output "my_cluster_client_key" {
  value = ovh_cloud_project_kube.my_cluster.kubeconfig_attributes[0].client_key
  sensitive = true
}
