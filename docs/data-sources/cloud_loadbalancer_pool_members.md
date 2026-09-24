---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer_pool_members (Data Source)

Use this data source to list the members of a pool in a public cloud load balancer.

## Example Usage

```terraform
data "ovh_cloud_loadbalancer_pool_members" "members" {
  service_name    = "<public cloud project ID>"
  loadbalancer_id = "<load balancer ID>"
  pool_id         = "<pool ID>"
}

output "member_addresses" {
  value = [for m in data.ovh_cloud_loadbalancer_pool_members.members.members : m.address]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `loadbalancer_id` - (Required) ID of the load balancer.
* `pool_id` - (Required) ID of the pool.

## Attributes Reference

The following attributes are exported:

* `members` - List of members. Each element exports:
  * `id` - Member ID.
  * `name` - Member name.
  * `address` - IP address of the member.
  * `protocol_port` - Port used by the member to receive traffic.
  * `subnet_id` - ID of the subnet the member is in.
  * `weight` - Weight of the member in the pool (0-256).
  * `backup` - Whether the member only receives traffic when all non-backup members are down.
  * `monitor` - Health monitor address and port override for this member:
    * `address` - IP address used by the health monitor for this member.
    * `port` - Port used by the health monitor for this member.
  * `checksum` - Computed hash representing the current target specification value.
  * `created_at` - Creation date of the member.
  * `updated_at` - Last update date of the member.
  * `resource_status` - Member readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
  * `current_state` - Current state of the member:
    * `name` - Member name.
    * `address` - IP address of the member.
    * `protocol_port` - Port used by the member.
    * `weight` - Weight of the member.
    * `subnet_id` - ID of the subnet the member is in.
    * `operating_status` - Operating status of the member.
    * `provisioning_status` - Provisioning status of the member.
    * `backup` - Whether this member is a backup member.
    * `monitor` - Health monitor address and port override (same schema as above).
