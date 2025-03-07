---
subcategory : "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_project_kube_nodes (Data Source)

Use this data source to get a list of OVHcloud Managed Kubernetes nodes.

## Example Usage

```terraform
data "ovh_cloud_project_kube_nodes" "nodes" {
  service_name  = "XXXXXX"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxx"
}

output "nodes" {
  value = data.ovh_cloud_project_kube_nodes.nodes
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `kube_id` - The ID of the managed kubernetes cluster.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `kube_id` - See Argument Reference above.
* `nodes` - List of all nodes composing the kubernetes cluster
  * `created_at` - Creation date
  * `deployed_at` - (Optional) Date of the effective deployment
  * `flavor` - Flavor name
  * `id` - ID of the node
  * `instance_id` - Openstack ID of the underlying VM of the node
  * `is_up_to_date` - Is the node in the target version of the cluster
  * `name` - Name of the node
  * `node_pool_id` - Managed kubernetes node pool ID
  * `project_id` - Public cloud project ID
  * `status` - Current status
  * `updated_at` - Last update date
  * `version` - Version in which the node is
