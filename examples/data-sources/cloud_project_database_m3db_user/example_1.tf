data "ovh_cloud_project_database_m3db_user" "m3db_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "m3db_user_group" {
  value = data.ovh_cloud_project_database_m3db_user.m3db_user.group
}
