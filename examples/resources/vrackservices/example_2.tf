# Vrack Services created by Terraform with a basic configuration

locals {
  region = "eu-west-lim"
}

data "ovh_me" "myaccount" {}

resource "ovh_vrackservices" "my-vrackservices" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
  plan = [
    {
      plan_code = "vrack-services"
      duration = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region_name"
          value = local.region
        }
      ]
    }
  ]
  target_spec = {
    subnets = [
      {
        az   = "eu-west-lim-a"
        cidr = "192.168.0.0/24"
      },
    ]
    services = {
      service_range = [
        {
          cidr              = "192.168.0.0/29"
          subnet_cidr       = "192.168.0.0/24"
          service_endpoints = []
        },
      ]
    }
  }
}