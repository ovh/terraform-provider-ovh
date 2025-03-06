---
subcategory : "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_project_kube_iprestrictions (Data Source)

Use this data source to get a OVHcloud Managed Kubernetes Service cluster IP restrictions.

## Example Usage

```terraform
data "ovh_cloud_project_kube_iprestrictions" "ip_restrictions" {
  service_name = "XXXXXX"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "ips" {
  value = data.ovh_cloud_project_kube_iprestrictions.ip_restrictions.ips
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `kube_id` - The id of the managed kubernetes cluster.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `kube_id` - See Argument Reference above.
* `ips` - The list of CIDRs that restricts the access to the API server.
