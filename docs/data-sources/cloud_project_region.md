---
subcategory : "Cloud Project"
---

# ovh_cloud_project_region (Data Source)

Use this data source to retrieve information about a region associated with a public cloud project. The region must be associated with the project.

## Example Usage

```terraform
data "ovh_cloud_project_region" "GRA1" {
  service_name = "XXXXXX"
  name         = "GRA1"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `name` - (Required) The name of the region associated with the public cloud project.

## Attributes Reference

`id` is set to the ID of the project concatenated with the name of the region. In addition, the following attributes are exported:

* `continent_code` - the code of the geographic continent the region is running. E.g.: EU for Europe, US for America...
* `datacenter_location` - The location code of the datacenter. E.g.: "GRA", meaning Gravelines, for region "GRA1"
* `services` - The list of public cloud services running within the region
  * `name` - the name of the public cloud service
  * `status` - the status of the service
