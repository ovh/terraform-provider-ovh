---
layout: "ovh"
page_title: "OVH: dedicated_nasha_partition_access"
sidebar_current: "docs-ovh-resource-dedicated-nasha-partition-access"
description: |-
  Provides a resource for managing access rights to partitions on HA-NAS services
---

# ovh_dedicated_nasha_partition_access

Provides a resource for managing access rights to partitions on HA-NAS services

## Example Usage

```
resource "ovh_dedicated_nasha_partition_access" "foo" {
  service_name = "zpool-12345"
  partition_name = "foo"
  ip = "123.123.123.123/32"
  protocol = "NFS"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your HA-NAS (it has to be ordered via OVH interface)
* `partition_name` - (Required) name of the partition
* `ip` - (Required) ip block in x.x.x.x/x format
* `type` - (Required) one of "readwrite", "readonly"

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `partition_name` - See Argument Reference above.
* `ip` - See Argument Reference above.
* `type` - See Argument Reference above.

## Import

HA-NAS partition access can be imported using the `{service_name}/{partition_name}/{ip}`, e.g.  
`$ terraform import ovh_dedicated_nasha_partition_access.foo zpool-12345/foo/123.123.123.123%2F32`
