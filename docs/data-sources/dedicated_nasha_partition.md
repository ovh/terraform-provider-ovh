---
subcategory : "NAS-HA"
---

# ovh_dedicated_nasha_partition (Data Source)

Use this data source to retrieve information about a dedicated HA-NAS partition.

## Example Usage

```terraform
data "ovh_dedicated_nasha_partition" "my_nas_ha_partition" {
  service_name = "zpool-12345"
  name         = "my-zpool-partition"
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your dedicated HA-NAS.
* `name` - (Required) The name of your dedicated HA-NAS partition.

## Attributes Reference

`id` is set with the name of the dedicated HA-NAS partition. In addition, the following attributes are exported:

* `size` - size of the partition in GB
* `protocol` - one of "NFS", "CIFS" or "NFS_CIFS"
* `description` - A brief description of the partition
* `capacity` - Percentage of partition space used in %
* `used_by_snapshots` - Percentage of partition space used by snapshots in %
