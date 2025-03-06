---
subcategory : "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_project_kube_iprestrictions

Apply IP restrictions to an OVHcloud Managed Kubernetes cluster.

## Example Usage

```terraform
resource "ovh_cloud_project_kube_iprestrictions" "vrack_only" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
  ips          = ["10.42.0.0/16"]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `kube_id` - The id of the managed Kubernetes cluster. **Changing this value recreates the resource.**
* `ips` - List of CIDR authorized to interact with the managed Kubernetes cluster.

## Attributes Reference

No additional attributes than the ones provided are exported.

## Timeouts

```terraform
resource "ovh_cloud_project_kube_iprestrictions" "vrack_only" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}
```
* `create` - (Default 10m)
* `update` - (Default 5m)
* `delete` - (Default 5m)

## Import

OVHcloud Managed Kubernetes Service cluster IP restrictions can be imported using the `service_name` and the `id` of the cluster, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_kube_iprestrictions.iprestrictions service_name/kube_id
```
