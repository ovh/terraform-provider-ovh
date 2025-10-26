# Configure the OVHcloud Provider
terraform {
  required_providers {
    ovh = {
      source = "ovh/ovh"
      version = ">= 2.7.0"
    }
  }
}

provider "ovh" {
}
