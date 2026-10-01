---
subcategory : "vRack"
---

# ovh_vrack_public_routing_priority

Create a new publicRoutingPriority for the vrack in a given region.

## Example Usage

Create a basic public routing priority

```terraform
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
```

Use the vRack datasource and configure a public routing priority

```terraform
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
```

Import an existing public routing priority

```terraform
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
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your vrack
* `region` - (Required) The region you want the public routing priority to be created on.
* `availability_zones` - (Required) A list of objects that allows to define on which AZ the Block should be active in priority, it is important to keep the objects sorted by ascending "priority" order.

## Attributes Reference

No additional attribute is exported.

## Import

The public routing priority can be imported using the vRack `service_name` (vRack identifier) and the publicRoutingPriority (Id), separated by "/" E.g.,

```bash
$ terraform import ovh_vrack_public_routing_priority.vrack_publicRoutingPriority "<service_name>/<publicRoutingPriorityId>"
```
