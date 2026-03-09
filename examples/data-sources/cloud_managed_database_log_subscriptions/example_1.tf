data "ovh_cloud_managed_database_log_subscriptions" "subscriptions" {
    service_name = "XXX"
    engine       = "YYY"
    cluster_id   = "ZZZ"
}

output "subscription_ids" {
  value = data.ovh_cloud_managed_database_log_subscriptions.subscriptions.subscription_ids
}
