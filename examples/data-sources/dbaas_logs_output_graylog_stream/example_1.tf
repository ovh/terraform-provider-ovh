
data "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "ldp-xx-xxxxx"
  title        = "my stream"
}
