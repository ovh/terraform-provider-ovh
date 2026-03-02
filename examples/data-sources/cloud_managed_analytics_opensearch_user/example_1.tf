data "ovh_cloud_managed_analytics_opensearch_user" "os_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "os_user_acls" {
  value = data.ovh_cloud_managed_analytics_opensearch_user.os_user.acls
}
