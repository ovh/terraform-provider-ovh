---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_capabilities (Data Source)

Use this data source to get information about capabilities of a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_database_capabilities" "capabilities" {
  service_name  = "XXX"
}

output "capabilities_engine_name" {
  value = tolist(data.ovh_cloud_project_database_capabilities.capabilities[*].engines)[0]
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

## Attributes Reference

The following attributes are exported:

`id` is set to `service_name` value. In addition, the following attributes are exported:

* `engines` - Database engines available.
  * `default_version` - Default version used for the engine.
  * `description` - Description of the engine.
  * `name` - Engine name.
  * `ssl_modes` - SSL modes for this engine.
  * `versions` - Versions available for this engine.
* `flavors` - Flavors available.
  * `core` - Flavor core number.
  * `memory` - Flavor ram size in GB.
  * `name` - Name of the flavor.
  * `storage` - Flavor disk size in GB.
* `options` - Options available.
  * `name` - Name of the option.
  * `type` - Type of the option.
* `plans` - Plans available.
  * `backup_retention` - Automatic backup retention duration.
  * `description` - Description of the plan.
  * `name` - Name of the plan.
* `service_name` - See Argument Reference above.
