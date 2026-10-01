---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer_pool_member

Creates a member in a pool of a public cloud load balancer.

## Example Usage

### Basic Member

```terraform
resource "ovh_cloud_loadbalancer_pool_member" "web1" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  pool_id         = ovh_cloud_loadbalancer_pool.pool.id
  name            = "web-1"
  address         = "10.0.0.11"
  protocol_port   = 8080
  subnet_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

### Weighted Backup Member with Monitor Override

```terraform
resource "ovh_cloud_loadbalancer_pool_member" "web_backup" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  pool_id         = ovh_cloud_loadbalancer_pool.pool.id
  name            = "web-backup"
  address         = "10.0.0.20"
  protocol_port   = 8080
  weight          = 10
  backup          = true

  monitor = {
    address = "10.0.0.20"
    port    = 9090
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `loadbalancer_id` - (Required) ID of the load balancer. **Changing this value recreates the resource.**
* `pool_id` - (Required) ID of the pool. **Changing this value recreates the resource.**
* `address` - (Required) IP address of the member. **Changing this value recreates the resource.**
* `protocol_port` - (Required) Port used by the member to receive traffic. **Changing this value recreates the resource.**
* `subnet_id` - (Optional) ID of the subnet the member is in. **Changing this value recreates the resource.**
* `name` - (Optional) Member name.
* `weight` - (Optional) Weight of the member in the pool (0-256). A higher weight receives more traffic. If omitted, the value assigned by the API is stored in the state.
* `backup` - (Optional) When `true`, the member is a backup member and only receives traffic when all non-backup members are down. If omitted, the value assigned by the API is stored in the state.
* `monitor` - (Optional) Health monitor address and port override for this member:
  * `address` - (Optional) IP address used by the health monitor for this member.
  * `port` - (Optional) Port used by the health monitor for this member.

## Attributes Reference

The following attributes are exported:

* `id` - Member ID.
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
  * `monitor` - Health monitor address and port override:
    * `address` - IP address used by the health monitor.
    * `port` - Port used by the health monitor.

## Import

A cloud load balancer pool member can be imported using the `service_name`, `loadbalancer_id`, `pool_id`, and `member_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_loadbalancer_pool_member.web1
  id = "<service_name>/<loadbalancer_id>/<pool_id>/<member_id>"
}
```

```bash
$ terraform import ovh_cloud_loadbalancer_pool_member.web1 service_name/loadbalancer_id/pool_id/member_id
```
