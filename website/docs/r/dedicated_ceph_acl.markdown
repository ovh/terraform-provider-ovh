---
layout: "ovh"
page_title: "OVH: dedicated_ceph_acl"
sidebar_current: "docs-ovh-resource-dedicated-ceph-acl"
description: |-
  Creates a new ACL for dedicated CEPH cluster.
---

# ovh_dedicated_ceph_acl

Add a new access ACL for the given network/mask.

## Example Usage

```hcl
data "ovh_dedicated_ceph" "my-ceph" {
  service_name = "94d423da-0e55-45f2-9812-836460a19939"
}

resource "ovh_dedicated_ceph_acl" "my-acl" {
  service_name = data.ovh_dedicated_ceph.my-ceph.id
  network      = "1.2.3.4"
  netmask      = "255.255.255.255"
  family       = "IPv4"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your dedicated CEPH
* `network` - (Required) The network IP to authorize 
* `netmask` - (Required) The network mask to apply
* `family` - (Required) IP family. `IPv4` or `IPv6`

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `network` - See Argument Reference above.
* `netmask` - See Argument Reference above.
* `family` - See Argument Reference above.
