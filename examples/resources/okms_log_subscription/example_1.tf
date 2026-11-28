resource "ovh_dbaas_logs_output_graylog_stream" "kms" {
  service_name = "ldp-xx-xxxxx"
  title        = "KMS secret audit logs"
  description  = "OVHcloud KMS REST API audit logs of the secrets"
}

resource "ovh_okms_log_subscription" "secret" {
  okms_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  kind      = "audit-secret"
  stream_id = ovh_dbaas_logs_output_graylog_stream.kms.stream_id
}
