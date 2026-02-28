data "ovh_cloud_managed_analytics" "opensearch" {
  service_name  = "XXX"
  engine        = "opensearch"
  id            = "ZZZ"
}

resource "ovh_cloud_managed_analytics_opensearch_user" "user" {
  service_name  = data.ovh_cloud_managed_analytics.opensearch.service_name
  cluster_id    = data.ovh_cloud_managed_analytics.opensearch.id
  acls {
    pattern    = "logs_*"
    permission = "read"
  }
  acls {
    pattern    = "data_*"
    permission = "deny"
  }
  name          = "johndoe"
}

output "user_password" {
  value     = ovh_cloud_managed_analytics_opensearch_user.user.password
  sensitive = true
}
