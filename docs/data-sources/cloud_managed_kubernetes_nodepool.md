---
subcategory: "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes_nodepool (Data Source)

Use this data source to retrieve a nodepool of a Managed Kubernetes cluster.

## Example Usage

```terraform
data "ovh_cloud_managed_kubernetes_nodepool" "pool" {
  service_name          = "<Public cloud project id>"
  managed_kubernetes_id = "<cluster id>"
  id                    = "<nodepool id>"
}

output "ready_nodes" {
  value = data.ovh_cloud_managed_kubernetes_nodepool.pool.current_state.status.ready_nodes
}
```

## Argument Reference

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `managed_kubernetes_id` - (Required) ID of the Managed Kubernetes cluster.
* `id` - (Required) Nodepool ID.

## Attributes Reference

Every attribute of the [`ovh_cloud_managed_kubernetes_nodepool`](../resources/cloud_managed_kubernetes_nodepool.md)
resource is exported, read-only: `name`, `desired_nodes`, `node_config`,
`availability_zone`, `billing_type`, `anti_affinity`, `autoscaler`, `template`,
`checksum`, `created_at`, `updated_at`, `resource_status` and `current_state`.
