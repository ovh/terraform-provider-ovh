---
subcategory : "Managed Kubernetes Service"
---

# ovh_cloud_project_kube_nodepool (Data Source)

Use this data source to get a OVHcloud Managed Kubernetes node pool.

## Example Usage

```hcl
data "ovh_cloud_project_kube_nodepool" "nodepool" {
  service_name  = "XXXXXX"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxx"
  name          = "xxxxxx"
}

output "max_nodes" {
  value = data.ovh_cloud_project_kube_nodepool.nodepool.max_nodes
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `kube_id` - The id of the managed kubernetes cluster.

* `name` - The name of the node pool.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `kube_id` - See Argument Reference above.
* `name` - (Optional) The name of the nodepool.
  Changing this value recreates the resource.
  Warning: "_" char is not allowed!
* `flavor_name` - a valid OVHcloud public cloud flavor ID in which the nodes will be started.
  Ex: "b2-7". Changing this value recreates the resource.
  You can find the list of flavor IDs: https://www.ovhcloud.com/fr/public-cloud/prices/
* `desired_nodes` - number of nodes to start.
* `max_nodes` - maximum number of nodes allowed in the pool.
  Setting `desired_nodes` over this value will raise an error.
* `min_nodes` - minimum number of nodes allowed in the pool.
  Setting `desired_nodes` under this value will raise an error.
* `monthly_billed` - (Optional) should the nodes be billed on a monthly basis. Default to `false`.
* `anti_affinity` - (Optional) should the pool use the anti-affinity feature. Default to `false`.
* `autoscale` - (Optional) Enable auto-scaling for the pool. Default to `false`.
* `autoscaling_scale_down_unneeded_time_seconds` - (Optional) scaleDownUnneededTimeSeconds autoscaling parameter
  How long a node should be unneeded before it is eligible for scale down
* `autoscaling_scale_down_unready_time_seconds` - (Optional) scaleDownUnreadyTimeSeconds autoscaling parameter
  How long an unready node should be unneeded before it is eligible for scale down
* `autoscaling_scale_down_utilization_threshold` - (Optional) scaleDownUtilizationThreshold autoscaling parameter
  Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down
* `available_nodes` - Number of nodes which are actually ready in the pool
* `created_at` - Creation date
* `current_nodes` - Number of nodes present in the pool
* `desired_nodes` - Number of nodes you desire in the pool
* `flavor` - Flavor name
* `project_id` - Project id
* `size_status` - Status describing the state between number of nodes wanted and available ones
* `status` - Current status
* `up_to_date_nodes` - Number of nodes with the latest version installed in the pool
* `updated_at` - Last update date
