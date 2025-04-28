resource "ovh_cloud_project_kube_nodepool" "node_pool_multi_zones" {
  service_name       = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  kube_id            = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name               = "my-pool-zone-a" //Warning: "_" char is not allowed!
  flavor_name        = "b3-8"
  desired_nodes      = 3
  availability_zones = ["eu-west-par-a"] //Currently, only one zone is supported
}
