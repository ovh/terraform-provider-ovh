---
subcategory: "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes_nodepool_nodes (Data Source)

Use this data source to list the nodes of a Managed Kubernetes cluster nodepool.

## Example Usage

```terraform
data "ovh_cloud_managed_kubernetes_nodepool_nodes" "nodes" {
  service_name          = "<Public cloud project id>"
  managed_kubernetes_id = "<cluster id>"
  nodepool_id           = "<nodepool id>"
}

output "instance_ids" {
  value = [for n in data.ovh_cloud_managed_kubernetes_nodepool_nodes.nodes.nodes : n.current_state.instance_id if n.current_state.instance_id != null]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `managed_kubernetes_id` - (Required) ID of the Managed Kubernetes cluster.
* `nodepool_id` - (Required) ID of the nodepool.

## Attributes Reference

The following attributes are exported:

* `nodes` - Nodes of the nodepool:
  * `id` - Node ID.
  * `created_at` - Node creation date.
  * `updated_at` - Last node modification date.
  * `resource_status` - Lifecycle status of the node (`CREATING`, `READY`, `UPDATING`, `OUT_OF_SYNC`, `SUSPENDED`, `ERROR`, `DELETING`, `UNKNOWN`).
  * `current_state` - Observed state of the node:
    * `name` - Node name.
    * `nodepool_id` - ID of the nodepool the node belongs to.
    * `flavor` - Flavor of the Public Cloud instance backing the node.
    * `instance_id` - ID of the Public Cloud instance backing the node. Absent until the instance is provisioned.
    * `version` - Kubernetes version installed on the node.
    * `up_to_date` - Whether the node runs the current nodepool specification.
    * `deployed_at` - Date the node joined the cluster. Absent until the node has joined.
