resource "ovh_cloud_managed_kubernetes_nodepool" "pool" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name          = "my-pool"
  flavor_name   = "b3-8"
  desired_nodes = 3
  template {
    metadata {
      annotations = {
        k1 = "v1"
        k2 = "v2"
      }
      finalizers = []
      labels = {
        k3 = "v3"
        k4 = "v4"
      }
    }
    spec {
      unschedulable = false
      taints = [
        {
          effect = "PreferNoSchedule"
          key    = "k"
          value  = "v"
        }
      ]
    }
  }
}
