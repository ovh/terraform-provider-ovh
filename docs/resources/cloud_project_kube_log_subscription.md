---
subcategory : "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_project_kube_log_subscription

Creates a log subscription for a Managed Kubernetes cluster associated with a public cloud project.

## Example Usage

Create a log subscription for a Managed Kubernetes cluster.

```terraform
data "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "ldp-xx-xxxxx"
  title        = "my stream"
}

data "ovh_cloud_project_kube" "cluster" {
  service_name = "XXXXXX"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

resource "ovh_cloud_project_kube_log_subscription" "sub" {
  service_name = data.ovh_cloud_project_kube.cluster.service_name
  kube_id      = data.ovh_cloud_project_kube.cluster.id
  stream_id    = data.ovh_dbaas_logs_output_graylog_stream.stream.stream_id
  kind         = "audit"
}
```

~> **Note:** The OVHcloud APIv6 token used by the provider must have `GET` access on the following URI in its token ACL: `/dbaas/logs/*/operation/*`

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `kube_id` - (Required, Forces new resource) The id of the managed kubernetes cluster.
* `kind` - (Required, Forces new resource) Log kind name of this subscription. Only `audit` is currently supported.
* `stream_id` - (Required, Forces new resource) Id of the target Log data platform stream.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `kube_id` - See Argument Reference above.
* `kind` - See Argument Reference above.
* `stream_id` - See Argument Reference above.
* `subscription_id` - ID of the log subscription.
* `created_at` - Creation date of the subscription.
* `updated_at` - Last update date of the subscription.
* `resource` - Information about the subscribed resource.
* `resource.name` - Name of the subscribed resource, where the logs come from.
* `resource.type` - Type of the subscribed resource, where the logs come from.

## Import

OVHcloud Managed Kubernetes Service cluster log subscription can be imported using the `service_name`, `kube_id` and `subscription_id` separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_kube_log_subscription.sub service_name/kube_id/subscription_id
```
