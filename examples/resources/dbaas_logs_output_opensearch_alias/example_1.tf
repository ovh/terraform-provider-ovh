resource "ovh_dbaas_logs_output_opensearch_alias" "alias" {
  service_name = "...."
  description  = "my opensearch alias"
  suffix = "alias"
}
