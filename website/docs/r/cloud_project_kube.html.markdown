---
layout: "ovh"
page_title: "OVH: cloud_project_kube"
sidebar_current: "docs-ovh-resource-cloud-project-kube-x"
description: |-
  Creates a kubernetes managed cluster in a public cloud project.
---

# ovh_cloud_project_kube

Creates a kubernetes managed cluster in a public cloud project.

## Example Usage

```hcl
resource "ovh_cloud_project_kube" "mykube" {
   service_name = "94d423da0e5545f29812836460a19939"
   name         = "my_kube_cluster"
   region       = "GRA7"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `name` - (Optional) The name of the kubernetes cluster.

* `region` - a valid OVH public cloud region ID in which the kubernetes
   cluster will be available. Ex.: "GRA1". Defaults to all public cloud regions.
   Changing this value recreates the resource.

* `version` - (Optional) kubernetes version to use.
   Changing this value recreates the resource. Defaults to latest available.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `name` - See Argument Reference above.
* `region` - See Argument Reference above.
* `version` - See Argument Reference above.
* `control_plane_is_up_to_date` - True if control-plane is up to date.
* `is_up_to_date` - True if all nodes and control-plane are up to date.
* `next_upgrade_versions` - Kubernetes versions available for upgrade.
* `nodes_url` - Cluster nodes URL.
* `status` - Cluster status. Should be normally set to 'READY'.
* `update_policy` - Cluster update policy. Choose between [ALWAYS_UPDATE,MINIMAL_DOWNTIME,NEVER_UPDATE]'.
* `url` - Management URL of your cluster.
* `kubeconfig` - The kubeconfig file. Use this file to connect to your kubernetes cluster.
