data "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "ldp-xx-xxxxx"
  title        = "my stream"
}

data "ovh_cloud_project_database" "db" {
  service_name  = "XXX"
  engine        = "YYY"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_log_subscription" "subscription" {
	service_name = data.ovh_cloud_project_database.db.service_name
	engine       = data.ovh_cloud_project_database.db.engine
	cluster_id   = data.ovh_cloud_project_database.db.id
	stream_id    = data.ovh_dbaas_logs_output_graylog_stream.stream.id
}
