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
   flavor        = "b2-7"
   desired_nodes = 3
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `kube_id` - The id of the managed kubernetes cluster.

* `name` - (Optional) The name of the nodepool.
   Changing this value recreates the resource.

* `flavor` - a valid OVH public cloud flavor ID in which the nodes will be start.
   cluster will be available. Ex.: "GRA1". Defaults to all public cloud regions.
   Changing this value recreates the resource.

* `desired_nodes` - (Optional) number of nodes to start.

* `max_nodes` - (Optional) maximum number of nodes allowde in the pool.
   Setting `desired_nodes` over this value will raise an error.

* `min_nodes` - (Optional) minimum number of nodes allowde in the pool.
   Setting `desired_nodes` under this value will raise an error.

* `monthly_billed` - (Optional) should the nodes be billed on a monthly basis.

* `anti_affinity` - (Optional) should the pool use the anti-affinity feature.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `kube_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `flavor` - See Argument Reference above.
* `desired_nodes` - See Argument Reference above.
* `max_nodes` - See Argument Reference above.
* `min_nodes` - See Argument Reference above.
* `monthly_billed` - See Argument Reference above.
* `anti_affinity` - See Argument Reference above.
* `status` - Nodepool status. Should be normally set to 'READY'.
