---
layout: "ovh"
page_title: "OVH: ovh_dedicated_nasha_partition_snapshot"
sidebar_current: "docs-ovh-resource-dedicated-nasha-partition-snapshot"
description: |-
  Provides a resource for managing snapshot to partitions on HA-NAS services
---

# ovh_dedicated_nasha_partition_snapshot

Provides a resource for managing **snapshot** to partitions on HA-NAS services

## Example Usage

```
resource "ovh_dedicated_nasha_partition_snapshot" "my-partition" {
  service_name = "zpool-12345"
  partition_name = "my-partition"
  type = "day-3"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your HA-NAS (it has to be ordered via OVHcloud interface)
* `partition_name` - (Required) name of the partition
* `type` - (Required) Snapshot interval, allowed : day-1, day-2, day-3, day-7, hour-1, hour-6

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `partition_name` - See Argument Reference above.
* `type` - See Argument Reference above.

## Import

HA-NAS partition snapshot can be imported using the `{service_name}/{partition_name}/{type}`, e.g.  
`$ terraform import ovh_dedicated_nasha_partition_snapshot.foo zpool-12345/my-partition/day-3`
