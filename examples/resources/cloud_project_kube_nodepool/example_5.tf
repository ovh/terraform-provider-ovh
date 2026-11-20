# Node pool reserved to MKS system workloads
resource "ovh_cloud_project_kube_nodepool" "system" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name          = "system"
  flavor_name   = "b3-8"
  desired_nodes = 2
  template {
    metadata {
      annotations = {}
      finalizers  = []
      labels      = {}
    }
    spec {
      unschedulable = false
      taints = [
        {
          effect = "NoSchedule"
          key    = "CriticalAddonsOnly"
          value  = "true"
        }
      ]
    }
  }
}

# Node pool for your own workloads
resource "ovh_cloud_project_kube_nodepool" "apps" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name          = "apps"
  flavor_name   = "b3-16"
  desired_nodes = 3
}
