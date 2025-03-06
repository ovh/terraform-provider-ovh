---
subcategory : "NAS-HA"
---

# ovh_dedicated_nasha_partition

Provides a resource for managing partitions on HA-NAS services

## Example Usage

```terraform
resource "ovh_dedicated_nasha_partition" "my_partition" {
  service_name = "zpool-12345"
  name = "my-partition"
  size = 20
  protocol = "NFS"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your HA-NAS (it has to be ordered via OVHcloud interface)
* `name` - (Required) name of the partition
* `size` - (Required) size of the partition in GB
* `protocol` - (Required) one of "NFS", "CIFS" or "NFS_CIFS"
* `description` - A brief description of the partition

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `name` - See Argument Reference above.
* `size` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `description` - See Argument Reference above.
* `capacity` - Percentage of partition space used in %
* `used_by_snapshots` - Percentage of partition space used by snapshots in %

## Import

HA-NAS can be imported using the `{service_name}/{name}`, e.g.
`$ terraform import ovh_dedicated_nasha_partition.my-partition zpool-12345/my-partition`
