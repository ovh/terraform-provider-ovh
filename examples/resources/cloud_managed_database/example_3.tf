resource "ovh_cloud_managed_database" "mongodb" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description   = "my-first-mongodb"
  engine        = "mongodb"
  version       = "5.0"
  plan          = "production"
  nodes {
    region      = "SBG"
    subnet_id   = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  nodes {
    region      = "SBG"
    subnet_id   = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  nodes {
    region      = "SBG"
    subnet_id   = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
    network_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  }
  flavor        = "db1-30"
}
