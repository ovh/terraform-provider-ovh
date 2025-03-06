resource "ovh_dbaas_logs_cluster" "ldp" {
  service_name     = "ldp-xx-xxxxx"
  cluster_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  archive_allowed_networks       = ["10.0.0.0/16"]
  direct_input_allowed_networks  = ["10.0.0.0/16"]
  query_allowed_networks         = ["10.0.0.0/16"]
}
