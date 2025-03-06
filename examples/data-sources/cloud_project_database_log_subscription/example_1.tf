data "ovh_cloud_project_database_log_subscription" "subscription" {
    service_name = "VVV"
    engine       = "XXX"
    cluster_id   = "YYY"
    id           = "ZZZ"
}

output "subscription_ldp_name" {
  value = data.ovh_cloud_project_database_log_subscription.subscription.ldp_service_name
}
