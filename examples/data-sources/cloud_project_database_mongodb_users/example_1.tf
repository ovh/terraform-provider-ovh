data "ovh_cloud_project_database_mongodb_users" "all" {
  service_name = "XXX"
  cluster_id   = "YYY"
}

output "mongo_user_names" {
  value = [for u in data.ovh_cloud_project_database_mongodb_users.all.users : u.name]
}
