---
subcategory: "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes

Creates a Managed Kubernetes cluster in a public cloud project.

The admin kubeconfig is not stored on this resource: read it from the
[`ovh_cloud_managed_kubernetes_kubeconfig`](../data-sources/cloud_managed_kubernetes_kubeconfig.md) data source.

## Example Usage

```terraform
resource "ovh_cloud_managed_kubernetes" "cluster" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  plan         = "FREE"
  name         = "my-cluster"
  version      = "1.34"

  cni     = { type = "CANAL" }
  network = { type = "PUBLIC" }
}
```

### Private cluster in a vRack

```terraform
resource "ovh_cloud_managed_kubernetes" "private" {
  service_name = "<Public cloud project id>"
  region       = "EU-WEST-PAR"
  plan         = "STANDARD"
  name         = "my-private-cluster"

  cni = {
    type = "CILIUM"
    cilium = {
      hubble = { enabled = true }
    }
  }

  network = {
    type = "VRACK"
    vrack = {
      private_network_id = "<private network id>"
      subnet_id          = "<subnet id>"
      default_egress     = { type = "MANAGED_GATEWAY" }
    }
  }

  api_server = {
    access = {
      authorized_ips = {
        rules = [{ cidr = "203.0.113.0/24" }]
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `region` - (Required) Region of the cluster. **Changing this value recreates the resource.**
* `plan` - (Required) Platform tier: `FREE` or `STANDARD`. **Changing this value recreates the resource.**
* `cni` - (Required) Container Network Interface plugin.
  * `type` - (Required) CNI variant: `CANAL` (FREE plan) or `CILIUM` (STANDARD plan). **Changing this value recreates the resource.**
  * `canal` - (Optional) Canal settings, set when `type` is `CANAL`. Defaults to the kube-proxy settings chosen by the API. Terraform will only perform drift detection of its value when present in a configuration.
    * `kube_proxy` - (Optional) kube-proxy configuration.
      * `mode` - (Required by the API) `IPTABLES` or `IPVS`. **Changing this value recreates the resource.**
      * `iptables` - (Optional) Set when `mode` is `IPTABLES`: `min_sync_period`, `sync_period` (ISO 8601 durations, for example `PT60S`).
      * `ipvs` - (Optional) Set when `mode` is `IPVS`: `scheduler` (`DESTINATION_HASHING`, `LEAST_CONNECTION`, `NEVER_QUEUE`, `ROUND_ROBIN`, `SHORTEST_EXPECTED_DELAY`, `SOURCE_HASHING`), `min_sync_period`, `sync_period`, `tcp_fin_timeout`, `tcp_timeout`, `udp_timeout` (ISO 8601 durations).
  * `cilium` - (Optional) Cilium settings, set when `type` is `CILIUM`. Defaults to the settings chosen by the API. Terraform will only perform drift detection of its value when present in a configuration.
    * `cluster_id` - (Optional) ClusterMesh cluster identifier, between 1 and 255.
    * `cluster_mesh` - (Optional) `enabled`, and `api_server` with `service_type` (`LoadBalancer` or `NodePort`) and `node_port`.
    * `hubble` - (Optional) `enabled`, `relay.enabled`, and `ui` with `enabled`, `backend_resources` and `frontend_resources` (each with `limits` and `requests` holding `cpu` and `memory`).
* `network` - (Required) Network topology of the nodes.
  * `type` - (Required) `PUBLIC` or `VRACK`. **Changing this value recreates the resource.**
  * `vrack` - (Optional) Private network settings, required when `type` is `VRACK`.
    * `private_network_id` - (Required by the API) Private network ID. **Changing this value recreates the resource.**
    * `subnet_id` - (Required by the API) Subnet ID of the nodes. **Changing this value recreates the resource.**
    * `load_balancer_subnet_id` - (Optional) Subnet ID of the load balancers. **Changing this value recreates the resource.**
    * `default_egress` - (Optional) Default egress route of the nodes. Defaults to `MANAGED_GATEWAY` on `STANDARD` and `NODE_PUBLIC_INTERFACE` on `FREE`. Terraform will only perform drift detection of its value when present in a configuration.
      * `type` - (Required by the API) `MANAGED_GATEWAY`, `NODE_PUBLIC_INTERFACE` or `USER_GATEWAY`.
      * `user_gateway` - (Optional) Required when `type` is `USER_GATEWAY`: `ip` (required by the API) of the gateway.
* `name` - (Optional) Cluster name. Generated when omitted.
* `version` - (Optional) Kubernetes version: a minor (`1.34`) resolves to its latest patch, or a full version (`1.34.1-2`). Changing it upgrades the cluster in place. A full version is kept as declared after automatic patch upgrades: see `current_state.deployed_version` for the running one. Defaults to the latest version. Terraform will only perform drift detection of its value when present in a configuration.
* `ip_allocation_policy` - (Optional) Internal addressing: `pods` and `services` IPv4 CIDR blocks. Defaults to the CIDR blocks chosen by the API. Terraform will only perform drift detection of its value when present in a configuration. **Changing a configured value recreates the resource.**
* `api_server` - (Optional) Control-plane exposure.
  * `access.authorized_ips.rules` - (Optional) List of `{ cidr = "..." }` allowed to reach the API server. The API defaults this to a list with `0.0.0.0/0`. Terraform will only perform drift detection of its value when present in a configuration.
  * `admission_plugins` - (Optional) Map of admission plugin name to `{ enabled = true|false }`. Plugins the API enables by default are all listed under `current_state.api_server.admission_plugins`; only the declared ones are managed. Terraform will only perform drift detection of its value when present in a configuration.
* `authentication` - (Optional) Control-plane authentication.
  * `type` - (Required by the API) `OIDC`.
  * `oidc` - (Optional) OpenID Connect provider: `client_id` (required by the API), `issuer_url` (required by the API), `ca_content`, `groups_claim`, `groups_prefix`, `required_claim`, `signing_algs`, `username_claim`, `username_prefix`.
* `auto_upgrade` - (Optional) Automatic upgrade policy: `patch.enabled` (enabled by default). Terraform will only perform drift detection of its value when present in a configuration.

An optional argument omitted from the configuration keeps the value the API applies:
its default, or the value found on import. Terraform only performs drift detection
of an argument when it is present in the configuration, and removing an argument
from the configuration does not reset it: set it explicitly to change it.

## Attributes Reference

The following attributes are exported:

* `id` - Cluster ID.
* `checksum` - Computed hash of the target spec.
* `created_at` - Creation date.
* `updated_at` - Last update date.
* `resource_status` - Cluster readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`…).
* `current_state` - Observed state of the cluster. Same shape as the arguments above (`name`, `region`, `plan`, `version`, `api_server`, `authentication`, `auto_upgrade`, `cni`, `ip_allocation_policy`, `network`), plus:
  * `deployed_version` - Full Kubernetes version deployed, for example `1.34.1-2`.
  * `api_server.access.public_endpoint.url` - URL of the control plane, available once the cluster is ready.

## Import

A Managed Kubernetes cluster can be imported using the `service_name` and the
cluster `id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_managed_kubernetes.cluster
  id = "<service_name>/<cluster_id>"
}
```

```bash
$ terraform import ovh_cloud_managed_kubernetes.cluster service_name/cluster_id
```
