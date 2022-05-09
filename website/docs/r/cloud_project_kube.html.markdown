---
layout: "ovh"
page_title: "OVH: cloud_project_kube"
sidebar_current: "docs-ovh-resource-cloud-project-kube-x"
description: |-
  Creates a kubernetes managed cluster in a public cloud project.
---

# ovh_cloud_project_kube

Creates a OVH Managed Kubernetes Service cluster in a public cloud project.

## Example Usage

```hcl
resource "ovh_cloud_project_kube" "mykube" {
   service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
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

* `private_network_id` - (Optional) OpenStack private network (or vrack) ID to use.
   Changing this value recreates the resource. Defaults - not use private network.

## Attributes Reference

The following attributes are exported:

* `id` - Managed Kubernetes Service ID
* `service_name` - See Argument Reference above.
* `name` - See Argument Reference above.
* `region` - See Argument Reference above.
* `version` - See Argument Reference above.
* `private_network_id` - OpenStack private network (or vrack) ID to use.
* `control_plane_is_up_to_date` - True if control-plane is up to date.
* `is_up_to_date` - True if all nodes and control-plane are up to date.
* `next_upgrade_versions` - Kubernetes versions available for upgrade.
* `nodes_url` - Cluster nodes URL.
* `status` - Cluster status. Should be normally set to 'READY'.
* `update_policy` - Cluster update policy. Choose between [ALWAYS_UPDATE,MINIMAL_DOWNTIME,NEVER_UPDATE]'.
* `url` - Management URL of your cluster.
* `kubeconfig` - The kubeconfig file. Use this file to connect to your kubernetes cluster.

## Import

OVHcloud Managed Kubernetes Service clusters can be imported using the `serviceName` and the `id` of the cluster, separated by "/" E.g.,

```
$ terraform import ovh_cloud_project_kube.my_kube_cluster a6678gggjh76hggjh7f59/a123bc45-a1b2-34c5-678d-678ghg7676ebc
```
