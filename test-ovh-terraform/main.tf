terraform {
  required_providers {
    ovh = {
      source = "ovh/ovh"
    }
  }
  backend "local" {
    path = "/home/mpenicau/git/test-ovh-terraform/.backend"
  }
}

provider "ovh" {
  endpoint      = "ovh-eu"
    application_key = "11359f6ffa201d29"
    application_secret = "96990a9f80dd5145983acd714e8e217c"
    consumer_key = "ba38f3b03ed5147e6d0ca0c80454e48f"
}

data "ovh_cloud_project" "project" {
  service_name = "marcellin-terraform-pci"
}
