---
subcategory : "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes (Data Source)

Use this data source to get a OVHcloud Managed Kubernetes Service cluster.

## Example Usage

```terraform
data "ovh_cloud_managed_kubernetes" "my_kube_cluster" {
  service_name = "XXXXXX"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "version" {
  value = data.ovh_cloud_managed_kubernetes.my_kube_cluster.version
}

output "kubeconfig" {
  value = data.ovh_cloud_managed_kubernetes.my_kube_cluster.kubeconfig
  sensitive = true
}

output "kube_host" {
  value = data.ovh_cloud_managed_kubernetes.my_kube_cluster.kubeconfig_attributes[0].host
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `kube_id` - The id of the managed kubernetes cluster.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `kube_id` - See Argument Reference above.
* `name` - The name of the managed kubernetes cluster.
* `region` - The OVHcloud public cloud region ID of the managed kubernetes cluster.
* `version` - Kubernetes version of the managed kubernetes cluster.
* `plan` - Plan of the managed kubernetes cluster.
* `private_network_id` - OpenStack private network (or vrack) ID to use.
* `load_balancers_subnet_id` - Openstack private network (or vRack) ID to use for load balancers.
* `nodes_subnet_id` - Openstack private network (or vRack) ID to use for nodes.
* `control_plane_is_up_to_date` - True if control-plane is up-to-date.
* `is_up_to_date` - True if all nodes and control-plane are up-to-date.
* `next_upgrade_versions` - Kubernetes versions available for upgrade.
* `nodes_url` - Cluster nodes URL.
* `status` - Cluster status. Should be normally set to 'READY'.
* `update_policy` - Cluster update policy. Choose between [ALWAYS_UPDATE,MINIMAL_DOWNTIME,NEVER_UPDATE]'.
* `url` - Management URL of your cluster.
* `kubeconfig` - (Sensitive) Raw kubeconfig file content for connecting to the cluster.
* `kubeconfig_attributes` - (Sensitive) Structured kubeconfig data for connecting to the cluster.
  * `host` - Kubernetes API server endpoint.
  * `cluster_ca_certificate` - (Sensitive) Cluster certificate authority data.
  * `client_certificate` - (Sensitive) Client certificate data for authentication.
  * `client_key` - (Sensitive) Client private key data for authentication.
* `kube_proxy_mode` - Selected mode for kube-proxy.
* `customization` - **Deprecated** (Optional) Use `customization_apiserver` and `customization_kube_proxy` instead. Kubernetes cluster customization
  * `apiserver` - Kubernetes API server customization
  * `kube_proxy` - Kubernetes kube-proxy customization
* `customization_apiserver` - Kubernetes API server customization
  * `admissionplugins` - Kubernetes API server admission plugins customization
    * `enabled` - Array of admission plugins enabled, default is ["NodeRestriction","AlwaysPulImages"] and only these admission plugins can be enabled at this time.
    * `disabled` - Array of admission plugins disabled, default is [] and only AlwaysPulImages can be disabled at this time.
* `customization_kube_proxy` - Kubernetes kube-proxy customization
  * `iptables` - Kubernetes cluster kube-proxy customization of iptables specific config.
    * `sync_period` - Minimum period that iptables rules are refreshed, in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration format.
    * `min_sync_period` - Period that iptables rules are refreshed, in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration format.
  * `ipvs` - Kubernetes cluster kube-proxy customization of IPVS specific config (durations format is [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration.
    * `sync_period` - Minimum period that IPVS rules are refreshed, in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration format.
    * `min_sync_period` - Minimum period that IPVS rules are refreshed in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration.
    * `scheduler` - IPVS scheduler.
    * `tcp_timeout` - Timeout value used for idle IPVS TCP sessions in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration.
    * `tcp_fin_timeout` - Timeout value used for IPVS TCP sessions after receiving a FIN in RFC3339 duration.
    * `udp_timeout` - timeout value used for IPVS UDP packets in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) duration.
