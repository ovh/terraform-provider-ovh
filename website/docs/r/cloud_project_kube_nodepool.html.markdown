---
layout: "ovh"
page_title: "OVH: cloud_project_kube_nodepool"
sidebar_current: "docs-ovh-resource-cloud-project-kube-nodepool-x"
description: |-
  Creates a nodepool in a kubernetes managed cluster.
---

# ovh_cloud_project_kube_nodepool

Creates a nodepool in a kubernetes managed cluster.

## Example Usage

```hcl
resource "ovh_cloud_project_kube_nodepool" "pool" {
   service_name  = "94d423da0e5545f29812836460a19939"
   kube_id       = "9260267d-2bf9-4d9a-bb6e-24b8969f65e2 "
   name          = "my_pool"
   flavor_name   = "b2-7"
   desired_nodes = 3
   max_nodes     = 3
   min_nodes     = 3
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `kube_id` - The id of the managed kubernetes cluster.

* `name` - (Optional) The name of the nodepool.
   Changing this value recreates the resource.

* `flavor_name` - a valid OVH public cloud flavor ID in which the nodes will be start.
   cluster will be available. Ex.: "b2-7". Changing this value recreates the resource.

* `desired_nodes` - number of nodes to start.

* `max_nodes` - maximum number of nodes allowed in the pool.
   Setting `desired_nodes` over this value will raise an error.

* `min_nodes` - minimum number of nodes allowed in the pool.
   Setting `desired_nodes` under this value will raise an error.

* `monthly_billed` - (Optional) should the nodes be billed on a monthly basis. Default to `false`.

* `anti_affinity` - (Optional) should the pool use the anti-affinity feature. Default to `false`.

* `autoscale` - (Optional) Enable auto-scaling for the pool. Default to `false`.

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
* `up_to_date_nodes` - Number of nodes with latest version installed in the pool
* `updated_at` - Last update date
