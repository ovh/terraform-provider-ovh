data "ovh_dbaas_logs_input_engine" "logstash" {
  service_name = "ldp-xx-xxxxx"
  name          = "logstash"
  version       = "6.8"
  is_deprecated = true
}
