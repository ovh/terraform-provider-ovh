resource "ovh_cloud_project_kube" "my_cluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA7"
}
