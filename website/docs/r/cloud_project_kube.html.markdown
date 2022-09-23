---
layout: "ovh"
page_title: "OVH: cloud_project_kube"
sidebar_current: "docs-ovh-resource-cloud-project-kube-x"
description: |-
  Creates a kubernetes managed cluster in a public cloud project.
---

# ovh_cloud_project_kube

Creates a OVHcloud Managed Kubernetes Service cluster in a public cloud project.

## Example Usage

Simple Kubernetes cluster creation:

```hcl
resource "ovh_cloud_project_kube" "mycluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA7"
}
```

Kubernetes cluster creation with API Server AdmissionPlugins configuration:

```hcl
resource "ovh_cloud_project_kube" "mycluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	name          = "my_kube_cluster"
	region        = "GRA5"
	customization {
		apiserver {
			admissionplugins {
				enabled = ["NodeRestriction"]
				disabled = ["AlwaysPullImages"]
			}
		}
	}
}
```

Kubernetes cluster creation attached to a VRack in `GRA5` region:

```hcl
resource "ovh_vrack_cloudproject" "attach" {
	service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" # vrack ID
	project_id   = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" # Public Cloud service name
}

resource "ovh_cloud_project_network_private" "network" {
	service_name = ovh_vrack_cloudproject.attach.service_name
	vlan_id    = 0
	name       = "terraform_testacc_private_net"
	regions    = ["GRA5"]
	depends_on = [ovh_vrack_cloudproject.attach]
}

resource "ovh_cloud_project_network_private_subnet" "networksubnet" {
  service_name = ovh_cloud_project_network_private.network.service_name
  network_id   = ovh_cloud_project_network_private.network.id

  # whatever region, for test purpose
  region     = "GRA5"
  start      = "192.168.168.100"
  end        = "192.168.168.200"
  network    = "192.168.168.0/24"
  dhcp       = true
  no_gateway = false

  depends_on   = [ovh_cloud_project_network_private.network]
}

output "openstackID" {
    value = one(ovh_cloud_project_network_private.network.regions_attributes[*].openstackid)
}

resource "ovh_cloud_project_kube" "mycluster" {
	service_name  = var.service_name
	name          = "test-kube-attach"
	region        = "GRA5"

	private_network_id = tolist(ovh_cloud_project_network_private.network.regions_attributes[*].openstackid)[0]
   
	private_network_configuration {
		default_vrack_gateway              = ""
		private_network_routing_as_default = false
	}

	depends_on = [
		ovh_cloud_project_network_private.network
	]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `name` - (Optional) The name of the kubernetes cluster.

* `region` - a valid OVHcloud public cloud region ID in which the kubernetes
   cluster will be available. Ex.: "GRA1". Defaults to all public cloud regions.
   Changing this value recreates the resource.

* `version` - (Optional) kubernetes version to use.
   Changing this value updates the resource. Defaults to latest available.

* `customization` - (Optional) Customer customization object
  * apiserver - Kubernetes API server customization
    * admissionplugins - (Optional) Kubernetes API server admission plugins customization
        * enabled - (Optional) Array of admission plugins enabled, default is ["NodeRestriction","AlwaysPulImages"] and only these admission plugins can be enabled at this time. 
        * disabled - (Optional) Array of admission plugins disabled, default is [] and only AlwaysPulImages can be disabled at this time.

* `private_network_id` - (Optional) OpenStack private network ID to use.
   Changing this value delete the resource(including ETCD user data). Defaults - not use private network.
   
~> __WARNING__ Updating the private network ID resets the cluster so that all user data is deleted.

* `private_network_configuration` - (Optional) The private network configuration
  * default_vrack_gateway - If defined, all egress traffic will be routed towards this IP address, which should belong to the private network. Empty string means disabled.
  * private_network_routing_as_default - Defines whether routing should default to using the nodes' private interface, instead of their public interface. Default is false.

* `update_policy` - Cluster update policy. Choose between [ALWAYS_UPDATE, MINIMAL_DOWNTIME, NEVER_UPDATE].

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
* `url` - Management URL of your cluster.
* `kubeconfig` - The kubeconfig file. Use this file to connect to your kubernetes cluster.

## Import

OVHcloud Managed Kubernetes Service clusters can be imported using the `service_name` and the `id` of the cluster, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_kube.my_kube_cluster service_name/kube_id
```
