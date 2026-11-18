---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_capabilities (Data Source)

Use this data source to get what managed database services a public cloud project can create: engines and their versions, flavors, plans, regions, disk types and options, using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_database_capabilities" "capabilities" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}
```

## Argument Reference

* `service_name` - (Required) ID of the public cloud project.

## Attributes Reference

The following attributes are exported:

* `disks` - Disk types.
  * `iops_per_gigabyte` - IOPS per GB of disk. Null when the IOPS are fixed.
  * `maximum_iops` - Maximum IOPS per volume.
  * `maximum_throughput` - Maximum throughput in MB/s per volume. Null without limit.
  * `name` - Name of the disk type.
  * `throughput_per_gigabyte` - Throughput in MB/s per GB of disk. Null when it doesn't scale.
* `engines` - Engines.
  * `category` - Category of the engine: DATABASE or ANALYTICS.
  * `default_version` - Version used when the creation sets none.
  * `description` - Description of the engine.
  * `lifecycle` - Lifecycle: started_at, end_of_sale_at, end_of_life_at and status (BETA, STABLE, DEPRECATED, END_OF_SALE, END_OF_LIFE or UNAVAILABLE).
    * `end_of_life_at` - End of life date.
    * `end_of_sale_at` - End of sale date.
    * `started_at` - Release date.
    * `status` - Current status.
  * `name` - Name of the engine.
  * `order` - Display order.
  * `ssl_modes` - TLS modes of the connections.
  * `storage` - Storage strategy: DISTRIBUTED, NOT_APPLICABLE or REPLICATED.
  * `tags` - Display tags.
  * `versions` - Versions of the engine.
    * `default` - Whether it is the default version.
    * `lifecycle` - Lifecycle: started_at, end_of_sale_at, end_of_life_at and status (BETA, STABLE, DEPRECATED, END_OF_SALE, END_OF_LIFE or UNAVAILABLE).
      * `end_of_life_at` - End of life date.
      * `end_of_sale_at` - End of sale date.
      * `started_at` - Release date.
      * `status` - Current status.
    * `name` - Name of the version.
    * `tags` - Display tags.
* `flavors` - Flavors.
  * `generation` - Generation of the flavor: GEN1, GEN2 or GEN3.
  * `lifecycle` - Lifecycle: started_at, end_of_sale_at, end_of_life_at and status (BETA, STABLE, DEPRECATED, END_OF_SALE, END_OF_LIFE or UNAVAILABLE).
    * `end_of_life_at` - End of life date.
    * `end_of_sale_at` - End of sale date.
    * `started_at` - Release date.
    * `status` - Current status.
  * `name` - Name of the flavor.
  * `order` - Display order.
  * `specifications` - Specifications of the flavor.
    * `memory` - RAM size.
      * `unit` - Unit of the size: GB or MB.
      * `value` - Size, in the unit.
    * `vcpus` - Number of vCPU cores.
  * `tags` - Display tags.
* `options` - Options.
  * `name` - Name of the option.
  * `type` - Type of the option: BOOLEAN, DOUBLE, DURATION, LONG or STRING.
* `plans` - Plans.
  * `backup_retention` - Retention of the automatic backups, as a duration.
  * `description` - Description of the plan.
  * `lifecycle` - Lifecycle: started_at, end_of_sale_at, end_of_life_at and status (BETA, STABLE, DEPRECATED, END_OF_SALE, END_OF_LIFE or UNAVAILABLE).
    * `end_of_life_at` - End of life date.
    * `end_of_sale_at` - End of sale date.
    * `started_at` - Release date.
    * `status` - Current status.
  * `name` - Name of the plan.
  * `order` - Display order.
  * `tags` - Display tags.
* `regions` - Regions.
  * `lifecycle` - Lifecycle: started_at, end_of_sale_at, end_of_life_at and status (BETA, STABLE, DEPRECATED, END_OF_SALE, END_OF_LIFE or UNAVAILABLE).
    * `end_of_life_at` - End of life date.
    * `end_of_sale_at` - End of sale date.
    * `started_at` - Release date.
    * `status` - Current status.
  * `name` - Name of the region.
  * `order` - Display order.
  * `tags` - Display tags.
  * `type` - Type of the region: LOCAL-ZONE, REGION-1-AZ or REGION-3-AZ.
