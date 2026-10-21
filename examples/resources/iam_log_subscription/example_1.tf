resource "ovh_dbaas_logs_output_graylog_stream" "access_policy" {
  service_name = "ldp-xx-xxxxx"
  title        = "IAM access policy logs"
  description  = "OVHcloud IAM access policy logs"
}

resource "ovh_iam_log_subscription" "access_policy" {
  kind      = "default"
  stream_id = ovh_dbaas_logs_output_graylog_stream.access_policy.stream_id
}
