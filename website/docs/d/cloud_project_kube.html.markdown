---
layout: "ovh"
page_title: "OVH: cloud_project_kube"
sidebar_current: "docs-ovh-datasource-cloud-project-kube-x"
description: |-
  Get information & status of a Kubernetes managed cluster in a public cloud project.
---

# ovh_cloud_project_kube (Data Source)

Use this data source to get a OVHcloud Managed Kubernetes Service cluster.

## Example Usage

```hcl
data "ovh_cloud_project_kube" "my_kube_cluster" {
  service_name = "XXXXXX"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "version" {
  value = data.ovh_cloud_project_kube.my_kube_cluster.version
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `kube_id` - The id of the managed kubernetes cluster.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `kube_id` - See Argument Reference above.
* `name` - The name of the managed kubernetes cluster.
* `region` - The OVHcloud public cloud region ID of the managed kubernetes cluster.
* `version` - Kubernetes version of the managed kubernetes cluster.
* `private_network_id` - OpenStack private network (or vrack) ID to use.
* `control_plane_is_up_to_date` - True if control-plane is up to date.
* `is_up_to_date` - True if all nodes and control-plane are up to date.
* `next_upgrade_versions` - Kubernetes versions available for upgrade.
* `nodes_url` - Cluster nodes URL.
* `status` - Cluster status. Should be normally set to 'READY'.
* `update_policy` - Cluster update policy. Choose between [ALWAYS_UPDATE,MINIMAL_DOWNTIME,NEVER_UPDATE]'.
* `url` - Management URL of your cluster.
* `kube_proxy_mode` - Selected mode for kube-proxy.
* `customization` - Customer customization object
    * `apiserver` - Kubernetes API server customization
        * `admissionplugins` - Kubernetes API server admission plugins customization
            * `enabled` - Array of admission plugins enabled, default is ["NodeRestriction","AlwaysPulImages"] and only these admission plugins can be enabled at this time.
            * `disabled` - Array of admission plugins disabled, default is [] and only AlwaysPulImages can be disabled at this time.
  * `kube_proxy` - Kubernetes kube-proxy customization
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
