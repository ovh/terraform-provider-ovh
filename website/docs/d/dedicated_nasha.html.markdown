---
layout: "ovh"
page_title: "OVH: dedicated_nasha"
sidebar_current: "docs-ovh-datasource-dedicated-nasha"
description: |-
  Get information of a dedicated HA-NAS.
---

# ovh_dedicated_nasha (Data Source)

Use this data source to retrieve information about a dedicated HA-NAS.

## Example Usage

```hcl
data "ovh_dedicated_nasha" "foo" {
  service_name = "zpool-12345"
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your dedicated HA-NAS.

## Attributes Reference

`id` is set with the service_name of the dedicated HA-NAS.
In addition, the following attributes are exported:

* `service_name` - The storage service name
* `can_create_partition` - True, if partition creation is allowed on this HA-NAS
* `custom_name` - The name you give to the HA-NAS
* `datacenter` - area of HA-NAS
* `disk_type` - the disk type of the HA-NAS. Possible values are: `hdd`, `ssd`, `nvme`
* `ip` - Access IP of HA-NAS
* `monitored` - Send an email to customer if any issue is detected
* `zpool_capacity` - percentage of HA-NAS space used in %
* `zpool_size` - the size of the HA-NAS in GB