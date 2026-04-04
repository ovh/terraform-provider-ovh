# Example with 'terraform import' command

import {
  to = ovh_vrack_public_routing_priority.vrack_publicRoutingPriority
  id = "pn-000000/d468cb8-c34f-482d-b7db-2732d2154f1a"
}

# Variables
locals {
  region = "eu-west-par"
  vrack_name = "pn-000000"
}

# Always sort availabilityZone by ascending priority order
resource "ovh_vrack_public_routing_priority" "vrack_publicRoutingPriority" {
     service_name  = local.vrack_name
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
