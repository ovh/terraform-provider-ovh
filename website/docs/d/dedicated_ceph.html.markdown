---
layout: "ovh"
page_title: "OVH: cloud_region"
sidebar_current: "docs-ovh-datasource-dedicated-ceph-x"
description: |-
  Get information & status of a dedicated CEPH instance.
---

# ovh_dedicated_ceph (Data Source)

Use this data source to retrieve information about a dedicated CEPH. 

## Example Usage

```hcl
data "ovh_dedicated_ceph" "my-ceph" {
  service_name = "XXXXXX"
}
```

## Argument Reference


* `service_name` - (Required) The service name of the dedicated CEPH cluster.


## Attributes Reference

* `ceph_mons` - list of CEPH monitors IPs
* `ceph_version` - CEPH cluster version
* `crush_tunables` - CRUSH algorithm settings. Possible values
  * OPTIMAL
  * DEFAULT
  * LEGACY
  * BOBTAIL
  * ARGONAUT
  * FIREFLY
  * HAMMER
  * JEWEL
* `label` - CEPH cluster label
* `region` - cluster region
* `size` - Cluster size in TB
* `state` - the state of the cluster
* `status` - the status of the service
