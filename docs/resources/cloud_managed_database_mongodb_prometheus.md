---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_mongodb_prometheus

Creates the Prometheus endpoint, which it enables and whose credentials it manages, in a managed database service (mongodb), using the Cloud Databases API v2.

## Example Usage

```terraform
resource "ovh_cloud_managed_database_mongodb_prometheus" "prometheus" {
  service_name   = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  service_id     = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  password_reset = "2026-11-12T10:00:00Z"
}

output "prometheus_password" {
  value     = ovh_cloud_managed_database_mongodb_prometheus.prometheus.password
  sensitive = true
}
```

The API has no Prometheus object of its own: creating this resource turns on the Prometheus endpoint of the service, and destroying it turns the endpoint off. Leave `prometheus_enabled` unset on the service: both would change the same setting. The endpoint isn't available on every service: when the API answers 501 Not Implemented, the creation fails.

The API returns the password only at the creation and when `password_reset` resets it. An import leaves the password empty, and a `password_reset` in the configuration then resets it at the next apply.

A locked service refuses any change: destroying this resource then leaves the endpoint enabled, with a warning. Once the service is unlocked, set `prometheus_enabled` to `false` on the service to turn it off.

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used. Changing it recreates the Prometheus endpoint.
* `service_id` - (Required) ID of the service. Changing it recreates the Prometheus endpoint.
* `password_reset` - (Optional) Arbitrary string to change to reset the password. Removing it keeps the password.
* `timeouts` - (Optional) How long to wait for the service, as durations such as `1h30m`:
  * `create` - (Optional) Wait for the service to enable the endpoint. Defaults to `40m`.
  * `delete` - (Optional) Wait for the service to disable the endpoint. Defaults to `40m`.

## Attributes Reference

The following attributes are exported:

* `engine` - Engine of the service: mongodb.
* `id` - ID of the service.
* `password` - (Sensitive) Password for Prometheus authentication.
* `srv_domain` - SRV domain that resolves to the scraping targets.
* `username` - Username for Prometheus authentication.

## Timeouts

The creation and the deletion wait up to 40 minutes for the service to enable or disable the endpoint. The `timeouts` argument sets other durations.

## Import

It can be imported using its `service_name` and `service_id`, separated by "/".

```terraform
import {
  to = ovh_cloud_managed_database_mongodb_prometheus.prometheus
  id = "<service_name>/<service_id>"
}
```

```bash
$ terraform import ovh_cloud_managed_database_mongodb_prometheus.prometheus service_name/service_id
```
