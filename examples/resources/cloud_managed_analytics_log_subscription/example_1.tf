data "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "ldp-xx-xxxxx"
  title        = "my stream"
}

data "ovh_cloud_managed_analytics" "db" {
  service_name = "XXX"
  engine       = "YYY"
  id           = "ZZZ"
}

resource "ovh_cloud_managed_analytics_log_subscription" "subscription" {
  service_name = data.ovh_cloud_managed_analytics.db.service_name
  engine       = data.ovh_cloud_managed_analytics.db.engine
  cluster_id   = data.ovh_cloud_managed_analytics.db.id
  stream_id    = data.ovh_dbaas_logs_output_graylog_stream.stream.id
  kind         = "customer_logs"
}
