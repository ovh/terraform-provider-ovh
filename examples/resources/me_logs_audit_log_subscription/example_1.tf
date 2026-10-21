resource "ovh_dbaas_logs_output_graylog_stream" "audit" {
  service_name = "ldp-xx-xxxxx"
  title        = "account audit logs"
  description  = "OVHcloud account audit logs"
}

resource "ovh_me_logs_audit_log_subscription" "audit" {
  kind      = "default"
  stream_id = ovh_dbaas_logs_output_graylog_stream.audit.stream_id
}
