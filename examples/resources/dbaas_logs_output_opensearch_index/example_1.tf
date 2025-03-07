resource "ovh_dbaas_logs_output_opensearch_index" "index" {
  service_name = "...."
  description  = "my opensearch index"
  suffix = "index"
}
