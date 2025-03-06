resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "...."
  title        = "my stream"
  description  = "my graylog stream"
}
