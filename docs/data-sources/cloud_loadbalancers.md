---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancers (Data Source)

Use this data source to list the load balancers of a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_loadbalancers" "lbs" {
  service_name = "<public cloud project ID>"
}

output "loadbalancer_names" {
  value = [for lb in data.ovh_cloud_loadbalancers.lbs.loadbalancers : lb.name]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.

## Attributes Reference

The following attributes are exported:

* `loadbalancers` - List of load balancers. Each element exports:
  * `id` - Load balancer ID.
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
