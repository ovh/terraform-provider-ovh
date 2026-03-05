---
subcategory : "Gateway"
---

# ovh_cloud_gateway

Creates a gateway in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "xxxxxxxxxx"
  name         = "my-private-network"
  region       = "GRA1"
}

resource "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "my-subnet"
  cidr         = "10.0.0.0/24"
  region       = "GRA1"
}

resource "ovh_cloud_gateway" "gateway" {
  service_name             = ovh_cloud_network_private_vrack.network.service_name
  region                   = "GRA1"
  name                     = "my-gateway"
  model                    = "S"
  external_gateway_enabled = true

  subnets {
    id = ovh_cloud_network_private_vrack_subnet.subnet.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `region` - (Required) Region where the gateway will be created. **Changing this value recreates the resource.**
* `name` - (Required) Gateway name.
* `model` - (Required) Gateway sizing model (`S`, `M`, `L`, `XL`, `2XL`, `3XL`).
* `availability_zone` - (Optional) Availability zone for the gateway. **Changing this value recreates the resource.**
* `description` - (Optional) Gateway description.
* `external_gateway_enabled` - (Optional) Whether the external gateway is enabled.
* `subnets` - (Optional) Subnets to attach as router interfaces:
  * `id` - (Required) Subnet ID.

## Attributes Reference

The following attributes are exported:

* `id` - Gateway ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the gateway.
* `updated_at` - Last update date of the gateway.
* `resource_status` - Gateway readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the gateway:
  * `name` - Gateway name.
  * `description` - Gateway description.
  * `model` - Gateway sizing model.
  * `status` - OpenStack router status (`ACTIVE`, `BUILD`, `DOWN`, `ERROR`).
  * `external_gateway_enabled` - Whether the external gateway is enabled.
  * `external_ip` - External IP address assigned to the gateway.
  * `subnets` - Currently attached subnets:
    * `id` - Subnet ID.
  * `region` - Region.
  * `availability_zone` - Availability zone.

## Import

A cloud gateway can be imported using the `service_name` and `gateway_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_gateway.gateway
  id = "<service_name>/<gateway_id>"
}
```

```bash
$ terraform import ovh_cloud_gateway.gateway service_name/gateway_id
```
