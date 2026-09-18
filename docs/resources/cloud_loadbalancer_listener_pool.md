---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer_listener_pool

Creates a pool for a listener in a public cloud load balancer.

## Example Usage

### Basic Pool

```terraform
resource "ovh_cloud_loadbalancer_listener_pool" "pool" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  listener_id     = ovh_cloud_loadbalancer_listener.http.id
  name            = "my-pool"
  protocol        = "HTTP"
  algorithm       = "ROUND_ROBIN"
}
```

### Pool with Session Persistence

```terraform
resource "ovh_cloud_loadbalancer_listener_pool" "sticky" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  listener_id     = ovh_cloud_loadbalancer_listener.http.id
  name            = "sticky-pool"
  protocol        = "HTTP"
  algorithm       = "ROUND_ROBIN"

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "JSESSIONID"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `loadbalancer_id` - (Required) ID of the load balancer. **Changing this value recreates the resource.**
* `listener_id` - (Required) ID of the listener. **Changing this value recreates the resource.**
* `protocol` - (Required) Protocol used by the pool (`HTTP`, `HTTPS`, `PROXY`, `PROXYV2`, `SCTP`, `TCP`, `UDP`). **Changing this value recreates the resource.**
* `algorithm` - (Required) Load balancing algorithm (`LEAST_CONNECTIONS`, `ROUND_ROBIN`, `SOURCE_IP`, `SOURCE_IP_PORT`).
* `name` - (Optional) Pool name.
* `description` - (Optional) Pool description.
* `persistence` - (Optional) Session persistence configuration:
  * `type` - (Required) Session persistence type (`APP_COOKIE`, `HTTP_COOKIE`, `SOURCE_IP`).
  * `cookie_name` - (Optional) Cookie name for `APP_COOKIE` persistence type.

## Attributes Reference

The following attributes are exported:

* `id` - Pool ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the pool.
* `updated_at` - Last update date of the pool.
* `resource_status` - Pool readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the pool:
  * `name` - Pool name.
  * `description` - Pool description.
  * `protocol` - Protocol used by the pool.
  * `algorithm` - Load balancing algorithm.
  * `persistence` - Session persistence configuration:
    * `type` - Session persistence type.
    * `cookie_name` - Cookie name.
  * `operating_status` - Operating status of the pool.
  * `provisioning_status` - Provisioning status of the pool.
  * `region` - Region.
  * `availability_zone` - Availability zone.

## Import

A cloud load balancer listener pool can be imported using the `service_name`, `loadbalancer_id`, `listener_id`, and `pool_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_loadbalancer_listener_pool.pool
  id = "<service_name>/<loadbalancer_id>/<listener_id>/<pool_id>"
}
```

```bash
$ terraform import ovh_cloud_loadbalancer_listener_pool.pool service_name/loadbalancer_id/listener_id/pool_id
```
