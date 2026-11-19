data "ovh_cloud_managed_database_users" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "mysql"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
}
