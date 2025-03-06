data "ovh_cloud_project_database" "m3db" {
  service_name  = "XXX"
  engine        = "m3db"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_m3db_user" "user" {
  service_name  = data.ovh_cloud_project_database.m3db.service_name
  cluster_id    = data.ovh_cloud_project_database.m3db.id
  group         = "mygroup"
  name          = "johndoe"
}

output "user_password" {
  value     = ovh_cloud_project_database_m3db_user.user.password
  sensitive = true
}
