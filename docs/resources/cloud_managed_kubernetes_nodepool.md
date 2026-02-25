---
subcategory : "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes_nodepool

Creates a nodepool in a OVHcloud Managed Kubernetes Service cluster.

## Example Usage

Create a simple node pool in your Kubernetes cluster:

```terraform
resource "ovh_cloud_managed_kubernetes_nodepool" "node_pool" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name          = "my-pool-1" //Warning: "_" char is not allowed!
  flavor_name   = "b3-8"
  desired_nodes = 3
}
```

Create a node pool on a specific availability zones for Kubernetes cluster (with multi-zones support):

```terraform
resource "ovh_cloud_managed_kubernetes_nodepool" "node_pool_multi_zones" {
  service_name       = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  kube_id            = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name               = "my-pool-zone-a" //Warning: "_" char is not allowed!
  flavor_name        = "b3-8"
  desired_nodes      = 3
  availability_zones = ["eu-west-par-a"] //Currently, only one zone is supported
}
```

Create an advanced node pool in your Kubernetes cluster:

```terraform
resource "ovh_cloud_managed_kubernetes_nodepool" "pool" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  kube_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name          = "my-pool"
  flavor_name   = "b3-8"
  desired_nodes = 3
  template {
    metadata {
      annotations = {
        k1 = "v1"
        k2 = "v2"
      }
      finalizers = []
      labels = {
        k3 = "v3"
        k4 = "v4"
      }
    }
    spec {
      unschedulable = false
      taints = [
        {
          effect = "PreferNoSchedule"
          key    = "k"
          value  = "v"
        }
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `kube_id` - The id of the managed kubernetes cluster. **Changing this value recreates the resource.**
* `name` - (Optional) The name of the nodepool. Warning: `_` char is not allowed! **Changing this value recreates the resource.**
* `flavor_name` - a valid OVHcloud public cloud flavor ID in which the nodes will be started. Ex: "b2-7". You can find the list of flavor IDs: https://www.ovhcloud.com/fr/public-cloud/prices/. **Changing this value recreates the resource.**
* `desired_nodes` - number of nodes to start.
* `availability_zones` - (Optional) list of availability zones to associate the pool - **mandatory for multi-zone** cluster - only one zone is supported at the moment.
* `max_nodes` - maximum number of nodes allowed in the pool. Setting `desired_nodes` over this value will raise an error.
* `min_nodes` - minimum number of nodes allowed in the pool. Setting `desired_nodes` under this value will raise an error.
* `monthly_billed` - (Optional) should the nodes be billed on a monthly basis. Default to `false`. **Changing this value recreates the resource.**
* `anti_affinity` - (Optional) should the pool use the anti-affinity feature. Default to `false`. **Changing this value recreates the resource.**
* `autoscale` - (Optional) Enable auto-scaling for the pool. Default to `false`.
* `autoscaling_scale_down_unneeded_time_seconds` - (Optional) scaleDownUnneededTimeSeconds autoscaling parameter How long a node should be unneeded before it is eligible for scale down
* `autoscaling_scale_down_unready_time_seconds` - (Optional) scaleDownUnreadyTimeSeconds autoscaling parameter How long an unready node should be unneeded before it is eligible for scale down
* `autoscaling_scale_down_utilization_threshold` - (Optional) scaleDownUtilizationThreshold autoscaling parameter Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down
* `template ` - (Optional) Managed Kubernetes nodepool template, which is a complex object constituted by two main nested objects:
  * `metadata` - Metadata of each node in the pool
    * `annotations` - Annotations to apply to each node
    * `finalizers` - Finalizers to apply to each node. A finalizer name must be fully qualified, e.g. kubernetes.io/pv-protection , where you prefix it with hostname of your service which is related to the controller responsible for the finalizer.
    * `labels` - Labels to apply to each node
  * `spec` - Spec of each node in the pool
    * `taints` - Taints to apply to each node [NodeSpec kubernetes documentation](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/node-v1/#NodeSpec)
      * `effect` - mandatory possible values: NoExecute, NoSchedule, PreferNoSchedule
      * `key` - mandatory
      * `value` - (Optional)
    * `unschedulable` - If true, set nodes as un-schedulable

## Attributes Reference

In addition, the following attributes are exported:

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

## Timeouts

```terraform
resource "ovh_cloud_managed_kubernetes_nodepool" "pool" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}
```

* `create` - (Default 20m)
* `update` - (Default 10m)
* `delete` - (Default 10m)

## Import

OVHcloud Managed Kubernetes Service cluster node pool can be imported using the `service_name`, the `id` of the cluster, and the `id` of the nodepool separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_managed_kubernetes_nodepool.pool service_name/kube_id/poolid
```
