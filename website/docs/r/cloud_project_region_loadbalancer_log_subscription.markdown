---
subcategory : "Log Subscriptions"
---

# ovh_cloud_project_region_loadbalancer_log_subscription

Subscribe to a Managed Loadbalance Logs Service in a public cloud project.

## Example Usage

Create a subscription

```hcl
resource "ovh_cloud_project_region_loadbalancer_log_subscription" "subscription" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region_name = "yyyy"
  loadbalancer_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  kind = "haproxy"
  stream_id = "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `region_name` - A valid OVHcloud public cloud region name in which the loadbalancer will be available. Ex.: "GRA11". **Changing this value recreates the resource.**
* `loadbalancer_id` - Loadbalancer id to get the logs  **Changing this value recreates the resource.**
* `stream_id` - Data stream id to use for the subscription  **Changing this value recreates the resource.**
* `kind` - haproxy  **Changing this value recreates the resource.**

## Attributes Reference

The following attributes are exported:

* `service_name` - The id of the public cloud project.
* `region_name` - A valid OVHcloud public cloud region name in which the loadbalancer will be available.
* `loadbalancer_id` - Loadbalancer id to get the logs
* `stream_id` - Data stream id to use for the subscription
* `kind` - haproxy
* `created_at` - The date of the subscription creation
* `ldp_service_name` - LDP service name
* `operation_id` - The operation ID
* `resource_name` - The resource name
* `resource_type` - The resource type
* `updated_at` - The last update of the subscription
* `subscription_id` - The subscription id