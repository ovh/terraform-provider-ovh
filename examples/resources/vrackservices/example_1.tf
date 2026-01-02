# Vrack Service order

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
    subnets = []
  }
}
