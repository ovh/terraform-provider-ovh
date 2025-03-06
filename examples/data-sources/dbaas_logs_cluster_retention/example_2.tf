data "ovh_dbaas_logs_cluster_retention" "retention" {
  service_name = "ldp-xx-xxxxx"
  cluster_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  duration     = "P14D"
}
