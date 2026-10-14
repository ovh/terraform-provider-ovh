resource "ovh_cloud_project_database" "postgresql" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description  = "my-first-postgresql"
  engine       = "postgresql"
  version      = "18"
  plan         = "production"
  nodes {
    region = "GRA"
  }
  nodes {
    region = "GRA"
  }
  flavor = "b3-8"
}
