# Configuration with a managed service

locals {
  region = "eu-west-lim"
  efs_name = "example-efs-service-name-000e75d3d4c1"
}

data "ovh_me" "myaccount" {}

data "ovh_storage_efs" "my-efs" {
  service_name = local.efs_name
}

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
        cidr         = "192.168.0.0/24"
        display_name = "my.subnet"
        service_range = {
          cidr = "192.168.0.0/29"
        }
        vlan = 30
        service_endpoints = [
          {
            managed_service_urn = data.ovh_storage_efs.my-efs.iam.urn
          }
        ]
      },
    ]
  }
}