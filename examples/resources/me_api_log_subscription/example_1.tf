resource "ovh_dbaas_logs_output_graylog_stream" "activity" {
  service_name = "ldp-xx-xxxxx"
  title        = "account activity logs"
  description  = "OVHcloud account activity logs"
}

resource "ovh_me_api_log_subscription" "activity" {
  kind      = "default"
  stream_id = ovh_dbaas_logs_output_graylog_stream.activity.stream_id
}
