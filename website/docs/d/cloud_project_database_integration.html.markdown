---
layout: "ovh"
page_title: "OVH: cloud_project_database_integration"
sidebar_current: "docs-ovh-datasource-cloud-project-database-integration"
description: |-
  Get information about an integration of a database cluster associated with a public cloud project.
---

# ovh_cloud_project_database_integration (Data Source)

Use this data source to get information about an integration of a database cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_integration" "integration" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
  id            = "UUU"
}

output "integration_type" {
  value = data.ovh_cloud_project_database_integration.integration.type
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required, Forces new resource) The engine of the database cluster you want to add. You can find the complete list of available engine in the [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).
All engines available exept `mongodb`

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `id` - (Required) Integration ID

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `destination_service_id` - ID of the destination service.
* `engine` - See Argument Reference above.
* `id` - See Argument Reference above.
* `parameters` - Parameters for the integration.
* `service_name` - See Argument Reference above.
* `source_service_id` - ID of the source service.
* `status` - Current status of the integration.
* `type` - Type of the integration.
