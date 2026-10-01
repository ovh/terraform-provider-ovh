---
subcategory : "Cloud Project"
---

# ovh_cloud_region (Data Source)

Use this data source to retrieve information about a single region of a public cloud project, using the OVHcloud API v2.

## Example Usage

```terraform
data "ovh_cloud_region" "region" {
  service_name = "<public cloud project ID>"
  name         = "GRA11"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `name`         - (Required) The name of the region (e.g. `GRA11`).

## Attributes Reference

The following attributes are exported:

* `status`             - Region status (`ENABLED`, `DISABLED` or `MAINTENANCE`).
* `continent`          - Continent code of the region.
* `country`            - Country code of the region.
* `datacenter_name`    - Display name of the datacenter hosting the region.
* `availability_zones` - Availability zones available in the region.
* `services`           - Available OpenStack services in the region.
