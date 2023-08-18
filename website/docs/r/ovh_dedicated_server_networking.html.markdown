---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_networking

Manage dedicated server networking interface on SCALE and HIGH-GRADE range.

!> The API route targeted by this resource are restricted to OVHCloud users (`Internal API`) with additional restrictions.


## Example Usage

The following example aims to bind all interfaces in a vRack

```hcl
locals {
  dedicated_server = "nsXXXXXX.ip-XX-XXX-XX.eu"
}

data "ovh_dedicated_server" "server" {
  service_name = local.dedicated_server
}

resource "ovh_dedicated_server_networking" "server"{
  service_name = local.dedicated_server
  interfaces {
      macs = sort(flatten(data.ovh_dedicated_server.server.vnis.*.nics)),
      type = "vrack"
  }
} 
```

The following example aims to attach the server to two different vRack.

```hcl
locals {
  dedicated_server = "nsXXXXXX.ip-XX-XXX-XX.eu"
}

data "ovh_dedicated_server" "server" {
  service_name = local.dedicated_server
}

resource "ovh_dedicated_server_networking" "server"{
  service_name = local.dedicated_server
  # Bond the two first interfaces in the first vrack
  interfaces {
      macs = slice(sort(flatten(data.ovh_dedicated_server.server.vnis.*.nics)), 0, 2)
      type = "vrack"
  }
  # Bond the two last interfaces in the last vrack
  interfaces {
      macs = slice(sort(flatten(data.ovh_dedicated_server.server.vnis.*.nics)), 2, 4)
      type = "vrack"
  }
} 
```

## Arguments Reference

The following arguments are required:

* `service_name` - (String) The service_name of your dedicated server. The full list of available dedicated servers can be found using the `ovh_dedicated_servers` datasource.
* `interfaces` - (Block List, Min: 1, Max: 2) Interface or interfaces aggregation.

### Nested Schema for interfaces

Required:

* `macs` (List of String) List of mac addresses to bind together.
* `type` (String) Type of bonding to create.

## Attributes Reference

The following attributes are exported:

* `description` - Operation description.
* `status` - status of the networking configuration (should be `active`).

## Import

A dedicated server networking configuration can be imported using the `service_name`.

```bash
$ terraform import ovh_dedicated_server_networking.server service_name
```
