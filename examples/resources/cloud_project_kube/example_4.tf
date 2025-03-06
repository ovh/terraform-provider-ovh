resource "ovh_cloud_project_kube" "my_cluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA7"
}

provider "helm" {
  kubernetes {
    host                    = ovh_cloud_project_kube.my_cluster.kubeconfig_attributes[0].host
    client_certificate      = base64decode(ovh_cloud_project_kube.my_cluster.kubeconfig_attributes[0].client_certificate)
    client_key              = base64decode(ovh_cloud_project_kube.my_cluster.kubeconfig_attributes[0].client_key)
    cluster_ca_certificate  = base64decode(ovh_cloud_project_kube.my_cluster.kubeconfig_attributes[0].cluster_ca_certificate)
  }
}

# Ready to use Helm provider
