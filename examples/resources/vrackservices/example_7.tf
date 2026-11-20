# BGP option ordered on an existing Vrack Services imported into Terraform

locals {
  vrack_services_id = "vrs-xxx-xxx-xxx-xxx"
}

import {
  to = ovh_vrackservices.my-vrackservices
  id = local.vrack_services_id
}

resource "ovh_vrackservices" "my-vrackservices" {
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
