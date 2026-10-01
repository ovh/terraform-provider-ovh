---
subcategory : "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_project_kube_log_subscription (Data Source)

Use this data source to get a log subscription for a Managed Kubernetes cluster.

## Example Usage

```terraform
data "ovh_cloud_project_kube_log_subscription" "sub" {
  service_name    = "XXXXXX"
  kube_id         = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subscription_id = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
}

output "resource-name" {
  value = data.ovh_cloud_project_kube_log_subscription.sub.resource.0.name
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `kube_id` - The id of the managed kubernetes cluster.
* `subscription_id` - The id of the log subscription.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `kube_id` - See Argument Reference above.
* `subscription_id` - See Argument Reference above.
* `kind` - Log kind name of this subscription. Only `audit` is currently supported.
* `stream_id` - Id of the target Log data platform stream.
* `created_at` - Creation date of the subscription.
* `updated_at` - Last update date of the subscription.
* `resource` - Information about the subscribed resource.
* `resource.name` - Name of the subscribed resource, where the logs come from.
* `resource.type` - Type of the subscribed resource, where the logs come from.
