---
layout: "ovh"
page_title: "OVH: cloud_project_kube_nodes"
sidebar_current: "docs-ovh-datasource-cloud-project-kube-nodes-x"
description: |-
Get information & status of a Kubernetes managed nodes in a public cloud project.
---

# ovh_cloud_project_kube_nodepool (Data Source)

Use this data source to get a list of OVHcloud Managed Kubernetes nodes.

## Example Usage

```hcl
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

* `service_name` - (Optional) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `kube_id` - The id of the managed kubernetes cluster.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `kube_id` - See Argument Reference above.
* `nodes` - List of all nodes composing the kubernetes cluster
  * `created_at` - Creation date
  * `deployed_at` - (Optional) Date of the effective deployment
  * `flavor` - Flavor name
  * `id` - Id of the node
  * `instance_id` - Openstack id of the underlying VM of the node
  * `is_up_to_date` - Is the node in the target version of the cluster
  * `name` - Name of the node
  * `node_pool_id` - Managed kubernetes node pool id
  * `project_id` - Public cloud project id 
  * `status` - Current status
  * `updated_at` - Last update date
  * `version` - Version in which the node is
