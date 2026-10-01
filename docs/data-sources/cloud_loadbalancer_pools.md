---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer_pools (Data Source)

Use this data source to list the pools of a load balancer in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_loadbalancer_pools" "pools" {
  service_name    = "<public cloud project ID>"
  loadbalancer_id = "<load balancer ID>"
}

output "pool_names" {
  value = [for p in data.ovh_cloud_loadbalancer_pools.pools.pools : p.name]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `loadbalancer_id` - (Required) ID of the load balancer.

## Attributes Reference

The following attributes are exported:

* `pools` - List of pools. Each element exports:
  * `id` - Pool ID.
  * `name` - Pool name.
  * `description` - Pool description.
  * `protocol` - Protocol used by the pool (`HTTP`, `HTTPS`, `PROXY`, `PROXYV2`, `SCTP`, `TCP`, `UDP`).
  * `algorithm` - Load balancing algorithm (`LEAST_CONNECTIONS`, `ROUND_ROBIN`, `SOURCE_IP`).
  * `persistence` - Session persistence configuration:
    * `type` - Session persistence type (`APP_COOKIE`, `HTTP_COOKIE`, `SOURCE_IP`).
    * `cookie_name` - Cookie name for `APP_COOKIE` persistence type.
  * `health_monitor` - Health monitor configuration:
    * `type` - Health monitor type (`HTTP`, `HTTPS`, `PING`, `TCP`, `UDP_CONNECT`, `SCTP`, `TLS_HELLO`).
    * `delay` - Seconds between health checks.
    * `timeout` - Seconds to wait for a health check response.
    * `max_retries` - Number of consecutive health check failures before marking member as unhealthy.
    * `max_retries_down` - Number of consecutive health check failures before marking member as `ERROR`.
    * `name` - Health monitor name.
    * `url_path` - URL path for HTTP/HTTPS health checks.
    * `http_method` - HTTP method for health checks.
    * `http_version` - HTTP version for health checks.
    * `expected_codes` - Expected HTTP response codes.
    * `domain_name` - Domain name for health check requests.
  * `checksum` - Computed hash representing the current target specification value.
  * `created_at` - Creation date of the pool.
  * `updated_at` - Last update date of the pool.
  * `resource_status` - Pool readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
  * `current_state` - Current state of the pool:
    * `name` - Pool name.
    * `description` - Pool description.
    * `protocol` - Protocol used by the pool.
    * `algorithm` - Load balancing algorithm.
    * `persistence` - Session persistence configuration (same schema as above).
    * `health_monitor` - Health monitor configuration, plus:
      * `id` - Health monitor ID.
      * `operating_status` - Operating status of the health monitor.
      * `provisioning_status` - Provisioning status of the health monitor.
    * `operating_status` - Operating status of the pool.
    * `provisioning_status` - Provisioning status of the pool.
