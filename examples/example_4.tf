terraform {
  required_providers {
    ovh = {
      source = "terraform.local/local/ovh"
      version = "0.0.1"
    }
  }
}

data "ovh_me" "me" {}

output "me" {
  value = data.ovh_me.me
}
