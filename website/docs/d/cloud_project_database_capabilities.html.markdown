---
layout: "ovh"
page_title: "OVH: cloud_project_database_capabilities"
sidebar_current: "docs-ovh-datasource-cloud-project-database-capabilities"
description: |-
  Get information about capabilities of a public cloud project.
---

# ovh_cloud_project_database_capabilities (Data Source)

Use this data source to get information about capabilities of a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_capabilities" "capabilities" {
  service_name  = "XXX"
}

output "capabilities_engine_name" {
  value = data.ovh_cloud_project_database_capabilities.capabilities.engine.0.name
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

## Attributes Reference

The following attributes are exported:

`id` is set to `service_name` value. In addition,
the following attributes are exported:

* `availability` - Availability of databases engines on cloud projects.
  * `backup` - Defines the type of backup.
  * `default` - Whether this availability can be used by default.
  * `end_of_life` - End of life of the product.
  * `engine` - Database engine name.
  * `flavor` - Flavor name.
  * `max_disk_size` - Maximum possible disk size in GB.
  * `max_node_number` - Maximum nodes of the cluster.
  * `min_disk_size` - Minimum possible disk size in GB.
  * `min_node_number` - Minimum nodes of the cluster.
  * `network` - Type of network.
  * `plan` - Plan name.
  * `region` - Region name.
  * `start_date` - Date of the release of the product.
  * `status` - Status of the availability.
  * `step_disk_size` - Flex disk size step in GB.
  * `upstream_end_of_life` - End of life of the upstream product.
  * `version` - Version name.
* `engine` - Database engines available.
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

