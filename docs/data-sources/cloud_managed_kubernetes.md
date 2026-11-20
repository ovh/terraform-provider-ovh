---
subcategory: "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes (Data Source)

Use this data source to retrieve a Managed Kubernetes cluster. Its admin kubeconfig is read with the
[`ovh_cloud_managed_kubernetes_kubeconfig`](cloud_managed_kubernetes_kubeconfig.md) data source.

## Example Usage

```terraform
data "ovh_cloud_managed_kubernetes" "cluster" {
  service_name = "<Public cloud project id>"
  id           = "<cluster id>"
}

output "api_server_url" {
  value = data.ovh_cloud_managed_kubernetes.cluster.current_state.api_server.access.public_endpoint.url
}
```

## Argument Reference

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `id` - (Required) Cluster ID.

## Attributes Reference

Every attribute of the [`ovh_cloud_managed_kubernetes`](../resources/cloud_managed_kubernetes.md)
resource is exported, read-only: `region`, `plan`, `name`, `version`, `cni`,
`network`, `ip_allocation_policy`, `api_server`, `authentication`,
`auto_upgrade`, `checksum`, `created_at`, `updated_at`, `resource_status` and
`current_state`.
