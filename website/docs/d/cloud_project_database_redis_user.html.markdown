---
layout: "ovh"
page_title: "OVH: cloud_project_database_redis_user"
sidebar_current: "docs-ovh-datasource-cloud-project-database-redis-user"
description: |-
  Get information about a user of a redis cluster associated with a public cloud project.
---

# ovh_cloud_project_database_redis_user (Data Source)

Use this data source to get information about a user of a redis cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_redis_user" "redisuser" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "redisuser_commands" {
  value = data.ovh_cloud_project_database_redis_user.redisuser.commands
}
```

## Argument Reference

* `service_name` - The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - Cluster ID

* `name` - Name of the user

## Attributes Reference

The following attributes are exported:

* `categories` - Categories of the user.
* `channels` - Channels of the user.
* `cluster_id` - See Argument Reference above.
* `commands` - Commands of the user.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `keys` - Keys of the user.
* `name` - See Argument Reference above.
* `service_name` - Current status of the user.
* `status` - Current status of the user.
