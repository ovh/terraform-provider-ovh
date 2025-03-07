data "ovh_cloud_project_database_mongodb_user" "mongo_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ@admin"
}

output "mongo_user_roles" {
  value = data.ovh_cloud_project_database_mongodb_user.mongo_user.roles
}
