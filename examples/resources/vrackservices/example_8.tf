# Vrack Services and BGP option ordered at once by Terraform

locals {
  region = "eu-west-lim"
  vrack_name = "pn-000000"
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
  plan_option = [
    {
      plan_code    = "bgp-beta-vrack-services"
      duration     = "P1M"
      pricing_mode = "default"
      quantity     = 1
    }
  ]
  target_spec = {
    subnets = [
      {
        az   = "eu-west-lim-a"
        cidr = "10.0.0.0/24"
      },
    ]
    bgp = {
      bgp_range = [
        {
          cidr                 = "10.0.0.0/28"
          subnet_cidr          = "10.0.0.0/24"
          client_peers         = ["10.0.0.5"]
          infrastructure_peers = ["10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4"]
        },
      ]
    }
  }
}

resource "ovh_vrack_vrackservices" "vrack-vrackservices-binding" {
  service_name   = local.vrack_name
  vrack_services = ovh_vrackservices.my-vrackservices.id
}
