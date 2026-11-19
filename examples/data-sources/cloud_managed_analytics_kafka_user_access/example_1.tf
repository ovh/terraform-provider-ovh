data "ovh_cloud_managed_analytics_kafka_user_access" "app" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  user_id      = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
}

output "app_key" {
  value     = data.ovh_cloud_managed_analytics_kafka_user_access.app.key
  sensitive = true
}
