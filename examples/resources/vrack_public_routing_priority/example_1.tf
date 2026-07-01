# Always sort availabilityZone by ascending priority order
resource "ovh_vrack_public_routing_priority" "vrack_publicRoutingPriority" {
     service_name  = "pn-000000"
     region        = "eu-west-par"
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
