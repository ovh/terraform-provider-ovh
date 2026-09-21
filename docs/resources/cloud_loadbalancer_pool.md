---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer_pool

Creates a pool in a public cloud load balancer. A pool is attached to a listener through the listener's `default_pool_id`, or targeted by an L7 policy with the `REDIRECT_TO_POOL` action.

## Example Usage

### Basic Pool

```terraform
resource "ovh_cloud_loadbalancer_pool" "pool" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  name            = "my-pool"
  protocol        = "HTTP"
  algorithm       = "ROUND_ROBIN"
}

resource "ovh_cloud_loadbalancer_listener" "http" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  name            = "http-listener"
  protocol        = "HTTP"
  protocol_port   = 80
  default_pool_id = ovh_cloud_loadbalancer_pool.pool.id
}
```

### Pool with Session Persistence and Health Monitor

```terraform
resource "ovh_cloud_loadbalancer_pool" "sticky" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  name            = "sticky-pool"
  protocol        = "HTTP"
  algorithm       = "ROUND_ROBIN"

  persistence = {
    type        = "APP_COOKIE"
    cookie_name = "JSESSIONID"
  }

  health_monitor = {
    type           = "HTTP"
    delay          = 5
    timeout        = 3
    max_retries    = 3
    url_path       = "/healthz"
    http_method    = "GET"
    expected_codes = "200"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `loadbalancer_id` - (Required) ID of the load balancer. **Changing this value recreates the resource.**
* `protocol` - (Required) Protocol used by the pool (`HTTP`, `HTTPS`, `PROXY`, `PROXYV2`, `SCTP`, `TCP`, `UDP`). **Changing this value recreates the resource.**
* `algorithm` - (Required) Load balancing algorithm (`LEAST_CONNECTIONS`, `ROUND_ROBIN`, `SOURCE_IP`). `SOURCE_IP_PORT` is not accepted: it is implemented by the OVN provider of Octavia only, which is not enabled on OVHcloud Public Cloud.
* `name` - (Optional) Pool name.
* `description` - (Optional) Pool description.
* `persistence` - (Optional) Session persistence configuration:
  * `type` - (Required) Session persistence type (`APP_COOKIE`, `HTTP_COOKIE`, `SOURCE_IP`).
  * `cookie_name` - (Optional) Cookie name for `APP_COOKIE` persistence type.
* `health_monitor` - (Optional) Health monitor configuration:
  * `type` - (Required) Health monitor type (`HTTP`, `HTTPS`, `PING`, `TCP`, `UDP_CONNECT`, `SCTP`, `TLS_HELLO`). **Changing this value recreates the resource.**
  * `delay` - (Required) Seconds between health checks.
  * `timeout` - (Required) Seconds to wait for a health check response.
  * `max_retries` - (Required) Number of consecutive health check failures before marking a member as unhealthy (1-10).
  * `max_retries_down` - (Optional) Number of consecutive health check failures before marking a member as `ERROR` (1-10).
  * `name` - (Optional) Health monitor name.
  * `url_path` - (Optional) URL path for HTTP/HTTPS health checks.
  * `http_method` - (Optional) HTTP method for health checks (`GET`, `HEAD`, `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS`, `TRACE`).
  * `http_version` - (Optional) HTTP version for health checks (`1.0` or `1.1`).
  * `expected_codes` - (Optional) Expected HTTP response codes (e.g. `200`, `200-202`).
  * `domain_name` - (Optional) Domain name for health check requests.

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
  * `health_monitor` - Health monitor configuration (same schema as the `health_monitor` argument), plus:
    * `id` - Health monitor ID.
    * `operating_status` - Operating status of the health monitor.
    * `provisioning_status` - Provisioning status of the health monitor.
  * `operating_status` - Operating status of the pool.
  * `provisioning_status` - Provisioning status of the pool.

## Import

A cloud load balancer pool can be imported using the `service_name`, `loadbalancer_id`, and `pool_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_loadbalancer_pool.pool
  id = "<service_name>/<loadbalancer_id>/<pool_id>"
}
```

```bash
$ terraform import ovh_cloud_loadbalancer_pool.pool service_name/loadbalancer_id/pool_id
```
