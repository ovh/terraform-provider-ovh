---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer

Creates a load balancer in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_loadbalancer" "lb" {
  service_name = <Public cloud project id>
  name         = "my-loadbalancer"
  region       = "GRA1"
  flavor_name  = "SMALL"
  description  = "My load balancer"

  network = {
    id        = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    subnet_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }
}
```

Pin the VIP to a fixed address of the subnet:

```terraform
resource "ovh_cloud_loadbalancer" "lb" {
  service_name = <Public cloud project id>
  name         = "my-loadbalancer"
  region       = "GRA1"
  flavor_name  = "SMALL"

  network = {
    id        = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    subnet_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    ip        = "10.0.0.37"
  }
}
```

Associate an existing floating IP to the VIP:

```terraform
resource "ovh_cloud_loadbalancer" "lb" {
  service_name = <Public cloud project id>
  name         = "my-loadbalancer"
  region       = "GRA1"
  flavor_name  = "SMALL"

  network = {
    id        = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    subnet_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    ip        = "203.0.113.42"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `region` - (Required) Region where the load balancer will be created. **Changing this value recreates the resource.**
* `name` - (Required) Load balancer name.
* `network` - (Required) Network of the VIP. **Changing any value of this block recreates the resource.**
  * `id` - (Required) ID of the network for the VIP.
  * `subnet_id` - (Required) ID of the subnet for the VIP. The subnet must belong to the network above.
  * `ip` - (Optional) IP of the VIP. When it belongs to the subnet CIDR, it pins the fixed VIP address: it must be inside the subnet allocation pool, must not be the subnet gateway IP and must not already be taken by another port. When it is outside the subnet CIDR, it must be an existing floating IP of the project in that region, not already associated with a port, and the subnet must be attached to a router with an external gateway; the floating IP is then associated to the VIP port. Left empty, the address is picked automatically inside the subnet.
* `flavor_name` - (Required) Name of the load balancer flavor. Must be one of `SMALL`, `MEDIUM`, `LARGE`, `XL`. **Changing this value recreates the resource.**
* `availability_zone` - (Optional) Availability zone for the load balancer. **Changing this value recreates the resource.**
* `description` - (Optional) Load balancer description.

## Attributes Reference

The following attributes are exported:

* `id` - Load balancer ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the load balancer.
* `updated_at` - Last update date of the load balancer.
* `resource_status` - Load balancer readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the load balancer:
  * `name` - Load balancer name.
  * `description` - Load balancer description.
  * `operating_status` - Operating status of the load balancer.
  * `provisioning_status` - Provisioning status of the load balancer.
  * `region` - Region.
  * `availability_zone` - Availability zone.
  * `network` - VIP network:
    * `id` - Network ID.
    * `subnet_id` - Subnet ID.
    * `addresses` - Addresses carried by the VIP port:
      * `ip` - IP address.
      * `type` - Address type (`FIXED`, `FLOATING`).
  * `flavor` - Load balancer flavor reference:
    * `id` - Flavor ID.

## Import

A cloud load balancer can be imported using the `service_name` and `loadbalancer_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_loadbalancer.lb
  id = "<service_name>/<loadbalancer_id>"
}
```

```bash
$ terraform import ovh_cloud_loadbalancer.lb service_name/loadbalancer_id
```
