resource "ovh_cloud_managed_database" "db" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "postgresql"
  description  = "my-postgresql"
  region       = "eu-west-gra"
  version      = "17"
  plan         = "business"
  flavor       = "b3-8"
}
