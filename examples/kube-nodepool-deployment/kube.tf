resource "ovh_cloud_project_kube" "my_kube_cluster" {
  count        = var.nb_cluster
  service_name = var.service_name
  name         = "my-cluster-${count.index}"
  region       = "GRA5"
}

resource "ovh_cloud_project_kube_nodepool" "node_pool" {
  count = var.nb_cluster
  service_name  = var.service_name
  kube_id       = ovh_cloud_project_kube.my_kube_cluster[count.index].id
  name          = "my-pool-${count.index}" //Warning: "_" char is not allowed!
  flavor_name   = "b2-7"
  desired_nodes = 3
  max_nodes     = 3
  min_nodes     = 3
}

resource "local_file" "kubeconfig" {
    count       = var.nb_cluster
    content     = ovh_cloud_project_kube.my_kube_cluster[count.index].kubeconfig
    filename = "my-kube-cluster-${count.index}.yml"
}