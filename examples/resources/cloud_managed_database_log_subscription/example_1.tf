data "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "ldp-xx-xxxxx"
  title        = "my-stream"
}

resource "ovh_cloud_managed_database_log_subscription" "logs" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "postgresql"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  stream_id    = data.ovh_dbaas_logs_output_graylog_stream.stream.id
  kind         = "customer_logs"
}
