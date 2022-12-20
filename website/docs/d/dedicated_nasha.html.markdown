---
layout: "ovh"
page_title: "OVH: dedicated_nasha"
sidebar_current: "docs-ovh-datasource-dedicated-nasha"
description: |-
  Get information of a dedicated nasha.
---

# ovh_dedicated_nasha (Data Source)

Use this data source to retrieve information about a dedicated nasha.

## Example Usage

```hcl
data "ovh_dedicated_nasha" "foo" {
  service_name = "zpool-12345"
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your dedicated NASHA.

## Attributes Reference

`id` is set with the service_name of the dedicated nasha.
In addition, the following attributes are exported:

* `service_name` - The storage service name
* `can_create_partition` - True, if partition creation is allowed on this nas HA
* `custom_name` - The name you give to the nas
* `datacenter` - area of nas
* `disk_type` - the disk type of the nasHa. Possible values are: `hdd`, `ssd`, `nvme`
* `ip` - Access ip of nas
* `monitored` - Send an email to customer if any issue is detected
* `zpool_capacity` - percentage of nas space used in %
* `zpool_size` - the size of the nas in Go