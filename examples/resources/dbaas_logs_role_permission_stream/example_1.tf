resource "ovh_dbaas_logs_role_permission_stream" "permission" {
  service_name     = "ldp-xx-xxxxx"

  role_id = ovh_dbaas_logs_role.ro.id
  stream_id = ovh_dbaas_logs_output_graylog_stream.mystream.stream_id
}
