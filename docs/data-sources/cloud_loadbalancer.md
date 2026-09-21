---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer (Data Source)

Use this data source to get information about a load balancer in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_loadbalancer" "lb" {
  service_name = "<public cloud project ID>"
  id           = "<load balancer ID>"
}

output "loadbalancer_vip" {
  value = data.ovh_cloud_loadbalancer.lb.current_state.network.addresses
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `id` - (Required) ID of the load balancer.

## Attributes Reference

The following attributes are exported:

* `name` - Load balancer name.
* `description` - Load balancer description.
* `region` - Region where the load balancer is located.
* `availability_zone` - Availability zone for the load balancer.
* `network` - Network of the VIP:
  * `id` - ID of the network for the VIP.
  * `subnet_id` - ID of the subnet for the VIP.
  * `ip` - IP requested for the VIP.
* `flavor_name` - Name of the load balancer flavor.
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
