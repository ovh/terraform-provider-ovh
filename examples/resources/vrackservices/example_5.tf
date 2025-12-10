# Complete Vrack Services configuration
# Once this plan executed, your ovh_vrackservices resource must be updated in your state using :
#   `terraform plan -refresh-only`
#   `terraform apply -refresh-only -auto-approve`

locals {
  region = "eu-west-lim"
  efs_name = "example-efs-service-name-000e75d3d4c1"
  vrack_name = "pn-000000"
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
          value = locals.region
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

resource "ovh_vrack_vrackservices" "vrack-vrackservices-binding" {
  service_name   = local.vrack_name
  vrack_services = ovh_vrackservices.my-vrackservices.id
}
