---
layout: "ovh"
page_title: "OVH: cloud_project_database_user"
sidebar_current: "docs-ovh-datasource-cloud-project-database-user"
description: |-
  Get information about a user of a database cluster associated with a public cloud project.
---

# cloud_project_database_user (Data Source)

Use this data source to get information about a user of a database cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_user" "user" {
  service_name  = "XXX"
  engine	      = "YYY"
  cluster_id    = "ZZZ"
  name          = "UUU"
}

output "user_name" {
  value = data.ovh_cloud_project_database_user.user.name
}
```

## Argument Reference

* `service_name` - The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - The engine of the database cluster you want user information. To get a full list of available engine visit :
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).\
Available engines for this resource (other have specific resource):
  * `cassandra`
  * `kafka`
  * `kafkaConnect`
  * `mysql`

* `cluster_id` - Cluster ID

* `name` - Name of the user.

## Attributes Reference

The following attributes are exported:

* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `status` - Current status of the user.
* `name` - Name of the user.
