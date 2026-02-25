resource "ovh_cloud_managed_kubernetes_iprestrictions" "vrack_only" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
  ips          = ["10.42.0.0/16"]
}
