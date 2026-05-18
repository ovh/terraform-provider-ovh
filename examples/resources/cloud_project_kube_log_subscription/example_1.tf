data "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "ldp-xx-xxxxx"
  title        = "my stream"
}

data "ovh_cloud_project_kube" "cluster" {
  service_name = "XXXXXX"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

resource "ovh_cloud_project_kube_log_subscription" "sub" {
  service_name = data.ovh_cloud_project_kube.cluster.service_name
  kube_id      = data.ovh_cloud_project_kube.cluster.id
  stream_id    = data.ovh_dbaas_logs_output_graylog_stream.stream.stream_id
  kind         = "audit"
}
