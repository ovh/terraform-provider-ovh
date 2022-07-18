---
layout: "ovh"
page_title: "OVH: cloud_project_database_users"
sidebar_current: "docs-ovh-datasource-cloud-project-database-users"
description: |-
  Get the list of users of a database cluster associated with a public cloud project.
---

# cloud_project_database_users (Data Source)

Use this data source to get the list of users of a database cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_users" "users" {
  service_name = "XXXX"
  engine	   = "YYYY"
  cluster_id   = "ZZZ"
}

output "user_ids" {
  value = data.ovh_cloud_project_database_users.users.user_ids
}
```

## Argument Reference

* `service_name` - The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - The engine of the database cluster you want to list users. To get a full list of available engine visit.
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).

* `cluster_id` - Cluster ID

## Attributes Reference

The following attributes are exported:

* `user_ids` - The list of users ids of the database cluster associated with the project.
