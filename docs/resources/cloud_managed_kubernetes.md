---
subcategory : "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes

Creates a OVHcloud Managed Kubernetes Service cluster in a public cloud project.

## Example Usage

Create a simple Kubernetes Free cluster in `GRA11` region:

```terraform
resource "ovh_cloud_managed_kubernetes" "my_cluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA11"
}

resource "ovh_cloud_managed_kubernetes_nodepool" "node_pool_1" {
  service_name  = ovh_cloud_managed_kubernetes.my_cluster.service_name
  kube_id       = ovh_cloud_managed_kubernetes.my_cluster.id
  name          = "my-pool-1"
  flavor_name   = "b3-8"
  desired_nodes = 3
}
```

Create a simple Kubernetes Free cluster in `GRA11` region and export its kubeconfig file:

```terraform
resource "ovh_cloud_managed_kubernetes" "my_cluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA11"
}

output "kubeconfig_file" {
  value     = ovh_cloud_managed_kubernetes.my_cluster.kubeconfig
  sensitive = true
}
```

Create a simple Kubernetes Free cluster in `GRA11` region and read kubeconfig attributes:

-> Sensitive attributes cannot be displayed using `terraform output` command. You need to specify the output's name: `terraform output my_cluster_host`.

```terraform
resource "ovh_cloud_managed_kubernetes" "my_cluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA11"
}

output "my_cluster_host" {
  value = ovh_cloud_managed_kubernetes.my_cluster.kubeconfig_attributes[0].host
  sensitive = true
}

output "my_cluster_cluster_ca_certificate" {
  value = ovh_cloud_managed_kubernetes.my_cluster.kubeconfig_attributes[0].cluster_ca_certificate
  sensitive = true
}

output "my_cluster_client_certificate" {
  value = ovh_cloud_managed_kubernetes.my_cluster.kubeconfig_attributes[0].client_certificate
  sensitive = true
}

output "my_cluster_client_key" {
  value = ovh_cloud_managed_kubernetes.my_cluster.kubeconfig_attributes[0].client_key
  sensitive = true
}
```

Create a simple Kubernetes Free cluster in `GRA11` region and use kubeconfig with [Helm provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs):

```terraform
resource "ovh_cloud_managed_kubernetes" "my_cluster" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name         = "my_kube_cluster"
  region       = "GRA11"
}

provider "helm" {
  kubernetes {
    host                    = ovh_cloud_managed_kubernetes.my_cluster.kubeconfig_attributes[0].host
    client_certificate      = base64decode(ovh_cloud_managed_kubernetes.my_cluster.kubeconfig_attributes[0].client_certificate)
    client_key              = base64decode(ovh_cloud_managed_kubernetes.my_cluster.kubeconfig_attributes[0].client_key)
    cluster_ca_certificate  = base64decode(ovh_cloud_managed_kubernetes.my_cluster.kubeconfig_attributes[0].cluster_ca_certificate)
  }
}

# Ready to use Helm provider
```

Create a Kubernetes Free cluster in `GRA11` region with API Server AdmissionPlugins configuration:

```terraform
resource "ovh_cloud_managed_kubernetes" "my_cluster" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name          = "my_kube_cluster"
  region        = "GRA11"
  customization_apiserver {
      admissionplugins {
        enabled = ["NodeRestriction"]
        disabled = ["AlwaysPullImages"]
      }
  }
}
```

Create a Kubernetes Free cluster in `GRA11` region with Kube proxy configuration, by specifying iptables or ipvs configurations:

```terraform
resource "ovh_cloud_managed_kubernetes" "my_cluster" {
  service_name    = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name            = "my_kube_cluster"
  region          = "GRA11"
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

Kubernetes Free cluster creation on a private network / subnet in `GRA11` region with a managed gateway:

```terraform
resource "ovh_cloud_project_network_private" "network" {
  service_name = var.service_name # Public Cloud service name
  vlan_id     = 42
  name       = "terraform_testacc_private_net"
  regions    = ["GRA11"]
}

resource "ovh_cloud_project_network_private_subnet" "subnet" {
  service_name = var.service_name
  network_id   = ovh_cloud_project_network_private.network.id

  # whatever region, for test purpose
  region     = "GRA11"
  start      = "192.168.168.100"
  end        = "192.168.168.200"
  network    = "192.168.168.0/24"
  dhcp       = true
  no_gateway = false
}

resource "ovh_cloud_project_gateway" "gateway" {
  service_name = var.service_name
  name       = "gateway"
  model      = "s"
  region     = "GRA11"
  network_id = tolist(ovh_cloud_project_network_private.network.regions_attributes[*].openstackid)[0]
  subnet_id  = ovh_cloud_project_network_private_subnet.subnet.id
}

resource "ovh_cloud_managed_kubernetes" "my_cluster" {
  service_name  = var.service_name
  name          = "test-kube-attach"
  region        = "GRA11"

  private_network_id = tolist(ovh_cloud_project_network_private.network.regions_attributes[*].openstackid)[0]
  nodes_subnet_id = ovh_cloud_project_network_private_subnet.subnet.id
  private_network_configuration {
      default_vrack_gateway              = ""
      private_network_routing_as_default = false
  }
}
```

Create a multi-zone Kubernetes Standard cluster (on 3 availability zones):

```terraform
resource "ovh_cloud_project_network_private" "network" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" # Public Cloud service name
  vlan_id      = 84
  name         = "terraform_mks_multiaz_private_net"
  regions      = ["EU-WEST-PAR"]
}

resource "ovh_cloud_project_network_private_subnet" "subnet" {
  service_name = ovh_cloud_project_network_private.network.service_name
  network_id   = ovh_cloud_project_network_private.network.id

  # whatever region, for test purpose
  region     = "EU-WEST-PAR"
  start      = "192.168.142.100"
  end        = "192.168.142.200"
  network    = "192.168.142.0/24"
  dhcp       = true
  no_gateway = false
}

resource "ovh_cloud_project_gateway" "gateway" {
  service_name = ovh_cloud_project_network_private.network.service_name
  name       = "gateway"
  model      = "s"
  region     = "EU-WEST-PAR"
  network_id = tolist(ovh_cloud_project_network_private.network.regions_attributes[*].openstackid)[0]
  subnet_id  = ovh_cloud_project_network_private_subnet.subnet.id
}

resource "ovh_cloud_managed_kubernetes" "my_multizone_cluster" {
  service_name  = ovh_cloud_project_network_private.network.service_name
  name          = "multi-zone-mks"
  region        = "EU-WEST-PAR"
  plan          = "standard"

  private_network_id = tolist(ovh_cloud_project_network_private.network.regions_attributes[*].openstackid)[0]
  nodes_subnet_id    = ovh_cloud_project_network_private_subnet.subnet.id

  depends_on    = [ ovh_cloud_project_gateway.gateway ] //Gateway is mandatory for multizones cluster
}

resource "ovh_cloud_managed_kubernetes_nodepool" "node_pool_multi_zones_a" {
  service_name       = ovh_cloud_project_network_private.network.service_name
  kube_id            = ovh_cloud_managed_kubernetes.my_multizone_cluster.id
  name               = "my-pool-zone-a" //Warning: "_" char is not allowed!
  flavor_name        = "b3-8"
  desired_nodes      = 3
  availability_zones = ["eu-west-par-a"] //Currently, only one zone is supported
}

resource "ovh_cloud_managed_kubernetes_nodepool" "node_pool_multi_zones_b" {
  service_name       = ovh_cloud_project_network_private.network.service_name
  kube_id            = ovh_cloud_managed_kubernetes.my_multizone_cluster.id
  name               = "my-pool-zone-b"
  flavor_name        = "b3-8"
  desired_nodes      = 3
  availability_zones = ["eu-west-par-b"]
}

resource "ovh_cloud_managed_kubernetes_nodepool" "node_pool_multi_zones_c" {
  service_name       = ovh_cloud_project_network_private.network.service_name
  kube_id            = ovh_cloud_managed_kubernetes.my_multizone_cluster.id
  name               = "my-pool-zone-c"
  flavor_name        = "b3-8"
  desired_nodes      = 3
  availability_zones = ["eu-west-par-c"]
}

output "kubeconfig_file_eu_west_par" {
  value     = ovh_cloud_managed_kubernetes.my_multizone_cluster.kubeconfig
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `name` - (Optional) The name of the kubernetes cluster.
* `region` - a valid OVHcloud public cloud region ID in which the kubernetes cluster will be available. Ex.: "GRA9". Defaults to all public cloud regions. **Changing this value recreates the resource.**
* `version` - (Optional) kubernetes version to use. Changing this value updates the resource. Defaults to the latest available.
* `plan` - (Optional) Plan of the MKS cluster `free` or `standard`. Default to `free`. Migration to another plan is not implemented yet.
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
* `private_network_id` - (Optional) Private network ID to use. **Changing this value recreates the resource, including ETCD user data.** Defaults - not use private network.

~> **WARNING** Updating the private network ID resets the cluster so that all user data is deleted.

* `nodes_subnet_id` - (Optional) Subnet ID to use for nodes, this subnet must belong to `private_network_id`. Default uses the first subnet belonging to the private network with id `private_network_id`. This attribute requires `private_network_id` to be defined. **Cannot be updated, it can only be used at cluster creation or reset.**
* `load_balancers_subnet_id` - (Optional) Subnet ID to use for Public Load Balancers, this subnet must belong to  `private_network_id`. Defaults to the same subnet as the nodes (see `nodes_subnet_id`). Requires `private_network_id` to be defined. See more network requirements in the [documentation](https://help.ovhcloud.com/csm/fr-public-cloud-kubernetes-expose-applications-using-load-balancer?id=kb_article_view&sysparm_article=KB0062873) for more information.

* `private_network_configuration` - (Optional) The private network configuration. If this is set then the 2 parameters below shall be defined.
  * `default_vrack_gateway` - If defined, all egress traffic will be routed towards this IP address, which should belong to the private network. Empty string means disabled.
  * `private_network_routing_as_default` - Defines whether routing should default to using the nodes' private interface, instead of their public interface. Default is false.

  In order to use the gateway IP advertised by the private network subnet DHCP, the following configuration shall be used.

  ```terraform
  private_network_configuration {
      default_vrack_gateway              = ""
      private_network_routing_as_default = true
  }
  ```
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

```terraform
resource "ovh_cloud_managed_kubernetes" "my_kube_cluster" {
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
$ terraform import ovh_cloud_managed_kubernetes.my_kube_cluster service_name/kube_id
```
