# Example with vRack datasource

# Variables
locals {
  region = "eu-west-par"
  vrack_name = "pn-000000"
}

# Datasources
data "ovh_vrack" "my_vrack" {
  service_name = local.vrack_name
}
output "my_vrack" {
  value = data.ovh_vrack.my_vrack.service_name
}

# Resources
# Always sort availabilityZone by ascending priority order
resource "ovh_vrack_public_routing_priority" "vrack_publicRoutingPriority" {
     service_name  = data.ovh_vrack.my_vrack.service_name
     region        = local.region
     availability_zones = [
      {
        priority = 1
        name = "eu-west-par-b"
      },
      {
        priority = 2
        name = "eu-west-par-c"
      },
      {
        priority = 3
        name = "eu-west-par-a"
      }
    ]
}
