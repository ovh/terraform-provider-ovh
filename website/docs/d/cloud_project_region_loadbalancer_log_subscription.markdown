---
subcategory : "Managed Loadbalancer Log Subscription"
---

# ovh_cloud_project_region_loadbalancer_log_subscription

Get information about a subscription to a Managed Loadbalancer Logs Service in a public cloud project.

## Example Usage


```hcl
data "ovh_cloud_project_region_loadbalancer_log_subscription" "sub" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region_name = "gggg"
  loadbalancer_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  subscription_id = "zzzzzzzz-yyyy-xxxx-wwww-vvvvvvvvvvvv"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `region_name` - A valid OVHcloud public cloud region name in which the loadbalancer is available. Ex.: "GRA11". 
* `loadbalancer_id` - Loadbalancer id to get the logs
* `subscription` - Subscription id

## Attributes Reference

The following attributes are exported:

* `service_name` - The id of the public cloud project.
* `region_name` - A valid OVHcloud public cloud region name in which the loadbalancer will be available. Ex.: "GRA11". 
* `loadbalancer_id` - Loadbalancer id to get the logs
* `stream_id` - Data stream id to use for the subscription
* `kind` - Router used for forwarding log 
* `created_at` - The date of the subscription creation
* `ldp_service_name` - LDP service name
* `operation_id` - The operation ID
* `resource_name` - The resource name
* `resource_type` - The resource type
* `updated_at` - The last update of the subscription
* `subscription_id` - The subscription id