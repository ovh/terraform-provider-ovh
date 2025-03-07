data "ovh_cloud_project_database_user" "user" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
  name          = "UUU"
}

output "user_name" {
  value = data.ovh_cloud_project_database_user.user.name
}
