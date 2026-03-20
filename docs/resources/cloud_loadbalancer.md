---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer

Creates a load balancer in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_loadbalancer" "lb" {
  service_name   = "xxxxxxxxxx"
  name           = "my-loadbalancer"
  region         = "GRA1"
  vip_network_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  vip_subnet_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  flavor_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  description    = "My load balancer"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `region` - (Required) Region where the load balancer will be created. **Changing this value recreates the resource.**
* `name` - (Required) Load balancer name.
* `vip_network_id` - (Required) ID of the network for the VIP. **Changing this value recreates the resource.**
* `vip_subnet_id` - (Required) ID of the subnet for the VIP. **Changing this value recreates the resource.**
* `flavor_id` - (Required) ID of the load balancer flavor. **Changing this value recreates the resource.**
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
  * `vip_address` - VIP address of the load balancer.
  * `operating_status` - Operating status of the load balancer.
  * `provisioning_status` - Provisioning status of the load balancer.
  * `region` - Region.
  * `availability_zone` - Availability zone.
  * `vip_network` - VIP network reference:
    * `id` - Network ID.
  * `vip_subnet` - VIP subnet reference:
    * `id` - Subnet ID.
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
