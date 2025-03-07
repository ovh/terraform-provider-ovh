---
subcategory : "Log Subscriptions"
---

# ovh_cloud_project_region_loadbalancer_log_subscriptions

Get information about subscriptions to a Managed Loadbalancer Logs Service in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_region_loadbalancer_log_subscriptions" "subs" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region_name = "gggg"
  loadbalancer_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `region_name` - A valid OVHcloud public cloud region name in which the loadbalancer is available. Ex.: "GRA11".
* `loadbalancer_id` - Loadbalancer id to get the logs
* `kind` - (Optional) currently only "haproxy" is available

## Attributes Reference

The following attributes are exported:

* `service_name` - The id of the public cloud project.
* `region_name` - A valid OVHcloud public cloud region name in which the loadbalancer will be available. Ex.: "GRA11".
* `loadbalancer_id` - Loadbalancer id to get the logs
* `kind` - Router used for forwarding log
* `subscription_ids` - The list of the subscription id
