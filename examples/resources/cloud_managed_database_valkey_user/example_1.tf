data "ovh_cloud_managed_database" "valkey" {
  service_name  = "XXXX"
  engine        = "valkey"
  id            = "ZZZZ"
}

resource "ovh_cloud_managed_database_valkey_user" "user" {
  service_name  = data.ovh_cloud_managed_database.valkey.service_name
  cluster_id    = data.ovh_cloud_managed_database.valkey.id
  categories    = ["+@set", "+@sortedset"]
  channels      = ["*"]
  commands      = ["+get", "-set"]
  keys          = ["data", "properties"]
  name          = "johndoe"
}

output "user_password" {
  value     = ovh_cloud_managed_database_valkey_user.user.password
  sensitive = true
}
