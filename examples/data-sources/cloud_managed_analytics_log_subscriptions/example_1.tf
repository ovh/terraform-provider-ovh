data "ovh_cloud_managed_analytics_log_subscriptions" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "kafka"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  kind         = "customer_logs"
}
