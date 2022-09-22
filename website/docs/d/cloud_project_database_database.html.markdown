---
layout: "ovh"
page_title: "OVH: cloud_project_database_database"
sidebar_current: "docs-ovh-datasource-cloud-project-database-database"
description: |-
  Get information about a database of a database cluster associated with a public cloud project.
---

# ovh_cloud_project_database_database (Data Source)

Use this data source to get information about a database of a database cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_database" "database" {
  service_name  = "XXX"
  engine	      = "YYY"
  cluster_id    = "ZZZ"
  name          = "UUU"
}

output "database_name" {
  value = data.ovh_cloud_project_database_database.database.name
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required) The engine of the database cluster you want database information. To get a full list of available engine visit :
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).\
Available engines for this resource (other have specific resource):
  * `mysql`
  * `postgresql`

* `cluster_id` - (Required) Cluster ID

* `name` - (Required) Name of the database.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the database.
* `default` - Defines if the database has been created by default.
* `id` - ID of the database.
* `service_name` - Current status of the database.
* `name` - Name of the database.
