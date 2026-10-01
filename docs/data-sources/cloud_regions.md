---
subcategory : "Cloud Project"
---

# ovh_cloud_regions (Data Source)

Use this data source to list the regions available for a public cloud project, using the OVHcloud API v2. All pages are fetched, so the full list of regions is returned.

## Example Usage

```terraform
data "ovh_cloud_regions" "regions" {
  service_name = "<public cloud project ID>"
}

output "regions" {
  value = data.ovh_cloud_regions.regions.regions
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.

## Attributes Reference

The following attributes are exported:

* `regions` - The list of regions available for the project. Each region exports:
  * `name`               - Name of the region (e.g. `GRA11`).
  * `status`             - Region status (`ENABLED`, `DISABLED` or `MAINTENANCE`).
  * `continent`          - Continent code of the region.
  * `country`            - Country code of the region.
  * `datacenter_name`    - Display name of the datacenter hosting the region.
  * `availability_zones` - Availability zones available in the region.
  * `services`           - Available OpenStack services in the region.
