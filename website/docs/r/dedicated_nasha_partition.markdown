---
layout: "ovh"
page_title: "OVH: dedicated_nasha_partition"
sidebar_current: "docs-ovh-resource-dedicated-nasha-partition"
description: |-
  Provides a resource for managing partitions on NASHA services
---

# ovh_dedicated_nasha_partition

Provides a resource for managing partitions on NASHA services

## Example Usage

```
resource "ovh_dedicated_nasha_partition" "foo" {
  service_name = "zpool-12345"
  name = "foo"
  size = 20
  protocol = "NFS"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your NASHA (it has to be ordered via OVH interface)
* `name` - (Required) name of the partition
* `size` - (Required) size of the partition in GB
* `protocol` - (Required) one of "NFS", "CIFS" or "NFS_CIFS"

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `name` - See Argument Reference above.
* `size` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `capacity` - Reports stroage space useage in %.
* `used_by_snapshots` - Reports storage space used by snapshots in %.
