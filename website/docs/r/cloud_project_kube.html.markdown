---
subcategory : "Managed Kubernetes Service"
---

# ovh_cloud_project_kube

Creates a OVHcloud Managed Kubernetes Service cluster in a public cloud project.

## Example Usage

Create a simple Kubernetes cluster in `GRA7` region:

```hcl
resource "ovh_cloud_project_kube" "mycluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA7"
}
```

Create a simple Kubernetes cluster in `GRA7` region and export its kubeconfig file:

```hcl
resource "ovh_cloud_project_kube" "mycluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA7"
}

output "kubeconfig_file" {
  value     = ovh_cloud_project_kube.mycluster.kubeconfig
  sensitive = true
}
```

Create a simple Kubernetes cluster in `GRA7` region and read kubeconfig attributes:

-> Sensitive attributes cannot be displayed using `terraform output` command. You need to specify the output's name: `terraform output mycluster-host`.

```hcl
resource "ovh_cloud_project_kube" "mycluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA7"
}

output "mycluster-host" {
  value = ovh_cloud_project_kube.mycluster.kubeconfig_attributes[0].host
  sensitive = true
}

output "mycluster-cluster-ca-certificate" {
  value = ovh_cloud_project_kube.mycluster.kubeconfig_attributes[0].cluster_ca_certificate
  sensitive = true
}

output "mycluster-client-certificate" {
  value = ovh_cloud_project_kube.mycluster.kubeconfig_attributes[0].client_certificate
  sensitive = true
}

output "mycluster-client-key" {
  value = ovh_cloud_project_kube.mycluster.kubeconfig_attributes[0].client_key
  sensitive = true
}
```

Create a simple Kubernetes cluster in `GRA7` region and use kubeconfig with [Helm provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs):

```hcl
resource "ovh_cloud_project_kube" "mycluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA7"
}

provider "helm" {
  kubernetes {
    host                    = ovh_cloud_project_kube.mycluster.kubeconfig_attributes[0].host
    client_certificate      = base64decode(ovh_cloud_project_kube.mycluster.kubeconfig_attributes[0].client_certificate)
    client_key              = base64decode(ovh_cloud_project_kube.mycluster.kubeconfig_attributes[0].client_key)
    cluster_ca_certificate  = base64decode(ovh_cloud_project_kube.mycluster.kubeconfig_attributes[0].cluster_ca_certificate)
  }
}

# Ready to use Helm provider
```

Create a Kubernetes cluster in `GRA5` region with API Server AdmissionPlugins configuration:

```hcl
resource "ovh_cloud_project_kube" "mycluster" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name          = "my_kube_cluster"
  region        = "GRA5"
  customization_apiserver {
      admissionplugins {
        enabled = ["NodeRestriction"]
        disabled = ["AlwaysPullImages"]
      }
  }
}
```

Create a Kubernetes cluster in `GRA5` region with Kube proxy configuration, by specifying iptables or ipvs configurations:

```hcl
resource "ovh_cloud_project_kube" "mycluster" {
  service_name    = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name            = "my_kube_cluster"
  region          = "GRA5"
  kube_proxy_mode = "ipvs" # or "iptables"	
	
  customization_kube_proxy {
    iptables {
      min_sync_period = "PT0S"
      sync_period = "PT0S"
    }
        
    ipvs {
      min_sync_period = "PT0S"
      sync_period = "PT0S"
      scheduler = "rr"
      tcp_timeout = "PT0S"
      tcp_fin_timeout = "PT0S"
      udp_timeout = "PT0S"
    }
  }
}
```

Kubernetes cluster creation attached to a VRack in `GRA5` region with:

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

  depends_on = [ovh_cloud_project_network_private.network]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `name` - (Optional) The name of the kubernetes cluster.
* `region` - a valid OVHcloud public cloud region ID in which the kubernetes cluster will be available. Ex.: "GRA1". Defaults to all public cloud regions. **Changing this value recreates the resource.**
* `version` - (Optional) kubernetes version to use. Changing this value updates the resource. Defaults to the latest available.
* `kube_proxy_mode` - (Optional) Selected mode for kube-proxy. **Changing this value recreates the resource, including ETCD user data.** Defaults to `iptables`.
* `customization` - **Deprecated** (Optional) Use `customization_apiserver` and `customization_kube_proxy` instead. Kubernetes cluster customization
  * `apiserver` - Kubernetes API server customization
  * `kube_proxy` - Kubernetes kube-proxy customization
* `customization_apiserver` - Kubernetes API server customization
  * `admissionplugins` - (Optional) Kubernetes API server admission plugins customization
      * `enabled` - (Optional) Array of admission plugins enabled, default is ["NodeRestriction","AlwaysPulImages"] and only these admission plugins can be enabled at this time. 
      * `disabled` - (Optional) Array of admission plugins disabled, default is [] and only AlwaysPulImages can be disabled at this time.
* `customization_kube_proxy` - Kubernetes kube-proxy customization
  * `iptables` - (Optional) Kubernetes cluster kube-proxy customization of iptables specific config (durations format is RFC3339 duration, e.g. `PT60S`)
      * `sync_period` - (Optional) Minimum period that iptables rules are refreshed, in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration format (e.g. `PT60S`).
      * `min_sync_period` - (Optional) Period that iptables rules are refreshed, in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration format (e.g. `PT60S`). Must be greater than 0.
  * `ipvs` - (Optional) Kubernetes cluster kube-proxy customization of IPVS specific config (durations format is [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration, e.g. `PT60S`)
      * `sync_period` - (Optional) Minimum period that IPVS rules are refreshed, in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration format (e.g. `PT60S`).
      * `min_sync_period` - (Optional) Minimum period that IPVS rules are refreshed in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration (e.g. `PT60S`).
      * `scheduler` - (Optional) IPVS scheduler.
      * `tcp_timeout` - (Optional) Timeout value used for idle IPVS TCP sessions in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration (e.g. `PT60S`). The default value is `PT0S`, which preserves the current timeout value on the system.
      * `tcp_fin_timeout` - (Optional) Timeout value used for IPVS TCP sessions after receiving a FIN in RFC3339 duration (e.g. `PT60S`). The default value is `PT0S`, which preserves the current timeout value on the system.
      * `udp_timeout` - (Optional) timeout value used for IPVS UDP packets in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration (e.g. `PT60S`). The default value is `PT0S`, which preserves the current timeout value on the system.
* `private_network_id` - (Optional) OpenStack private network (or vRack) ID to use. **Changing this value recreates the resource, including ETCD user data.** Defaults - not use private network.
   
~> __WARNING__ Updating the private network ID resets the cluster so that all user data is deleted.

* `private_network_configuration` - (Optional) The private network configuration. If this is set then the 2 parameters below shall be defined.
  * `default_vrack_gateway` - If defined, all egress traffic will be routed towards this IP address, which should belong to the private network. Empty string means disabled. Hence if the DHCP service with a Gateway IP is set on the subnet, then this IP will be used.
  * `private_network_routing_as_default` - Defines whether routing should default to using the nodes' private interface, instead of their public interface. Default is false.
* `update_policy` - Cluster update policy. Choose between [ALWAYS_UPDATE, MINIMAL_DOWNTIME, NEVER_UPDATE].

## Attributes Reference

The following attributes are exported:

* `control_plane_is_up_to_date` - True if control-plane is up-to-date.
* `id` - Managed Kubernetes Service ID
* `is_up_to_date` - True if all nodes and control-plane are up-to-date.
* `kubeconfig` - The kubeconfig file. Use this file to connect to your kubernetes cluster.
* `kubeconfig_attributes` - The kubeconfig file attributes.
  * `host` - The kubernetes API server URL.
  * `cluster_ca_certificate` - The kubernetes API server CA certificate.
  * `client_certificate` - The kubernetes API server client certificate.
  * `client_key` - The kubernetes API server client key.
* `name` - See Argument Reference above.
* `next_upgrade_versions` - Kubernetes versions available for upgrade.
* `nodes_url` - Cluster nodes URL.
* `private_network_configuration` - See Argument Reference above.
* `private_network_id` - See Argument Reference above.
* `region` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `status` - Cluster status. Should be normally set to 'READY'.
* `update_policy` - See Argument Reference above.
* `url` - Management URL of your cluster.
* `version` - See Argument Reference above.
* `customization_apiserver` - See Argument Reference above.
* `customization_kube_proxy` - See Argument Reference above.

## Timeouts

```hcl
resource "ovh_cloud_project_kube" "my_kube_cluster" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}
```
* `create` - (Default 10m)
* `update` - (Default 10m)
* `delete` - (Default 10m)

## Import

OVHcloud Managed Kubernetes Service clusters can be imported using the `service_name` and the `id` of the cluster, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_kube.my_kube_cluster service_name/kube_id
```
