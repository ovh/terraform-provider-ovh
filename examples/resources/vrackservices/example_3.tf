# Vrack Services associated to a vRack

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
  target_spec = {
    subnets = []
  }
}

resource "ovh_vrack_vrackservices" "vrack-vrackservices-binding" {
  service_name   = local.vrack_name
  vrack_services = ovh_vrackservices.my-vrackservices.id
}
