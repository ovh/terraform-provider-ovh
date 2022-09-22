---
layout: "ovh"
page_title: "OVH: cloud_project_database_m3db_user"
sidebar_current: "docs-ovh-datasource-cloud-project-database-m3db-user"
description: |-
  Get information about a user of a M3DB cluster associated with a public cloud project.
---

# ovh_cloud_project_database_m3db_user (Data Source)

Use this data source to get information about a user of a M3DB cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_m3db_user" "m3dbuser" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "m3dbuser_group" {
  value = data.ovh_cloud_project_database_m3db_user.m3dbuser.group
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `name` - (Required, Forces new resource) Name of the user.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `group` - See Argument Reference above.
* `name` - See Argument Reference above.
* `service_name` - Current status of the user.
* `status` - Current status of the user.