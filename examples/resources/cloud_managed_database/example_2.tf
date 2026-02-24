resource "ovh_cloud_managed_database" "postgresql" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-postgresql"
  engine        = "postgresql"
  version       = "14"
  plan          = "business"
  nodes {
    region  = "GRA"
  }
  nodes {
    region  = "GRA"
  }
  flavor        = "db1-15"
}
