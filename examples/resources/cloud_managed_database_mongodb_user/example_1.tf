data "ovh_cloud_managed_database" "mongodb" {
  service_name  = "XXX"
  engine        = "mongodb"
  id            = "ZZZ"
}

resource "ovh_cloud_managed_database_mongodb_user" "user" {
  service_name  = data.ovh_cloud_managed_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_managed_database.mongodb.id
  name          = "johndoe"
  roles         = ["backup@admin", "readAnyDatabase@admin"]
}

output "user_password" {
  value     = ovh_cloud_managed_database_mongodb_user.user.password
  sensitive = true
}
