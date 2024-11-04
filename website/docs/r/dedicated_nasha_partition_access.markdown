---
subcategory : "NAS-HA"
---

# ovh_dedicated_nasha_partition_access

Provides a resource for managing access rights to partitions on HA-NAS services

## Example Usage

```hcl
resource "ovh_dedicated_nasha_partition_access" "my_partition" {
  service_name    = "zpool-12345"
  partition_name  = "my-partition"
  ip              = "123.123.123.123/32"
  type            = "readwrite"
  acl_description = "Description of the ACL"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your HA-NAS (it has to be ordered via OVHcloud interface)
* `partition_name` - (Required) Name of the partition
* `ip` - (Required) IP block in x.x.x.x/x format
* `type` - (Required) One of "readwrite", "readonly"
* `acl_description` - A brief description of the acl

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `partition_name` - See Argument Reference above.
* `ip` - See Argument Reference above.
* `type` - See Argument Reference above.
* `acl_description` - See Argument Reference above.

## Import

HA-NAS partition access can be imported using the `{service_name}/{partition_name}/{ip}`, e.g.  
`$ terraform import ovh_dedicated_nasha_partition_access.my-partition zpool-12345/my-partition/123.123.123.123%2F32`
