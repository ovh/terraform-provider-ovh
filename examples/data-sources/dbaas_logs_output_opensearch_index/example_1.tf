
data "ovh_dbaas_logs_output_opensearch_index" "index" {
  service_name = "ldp-xx-xxxxx"
  name        = "index-name"
}
