---
layout: "ovh"
page_title: "OVH: cloud_project_database_redis_user"
sidebar_current: "docs-ovh-resource-cloud-project-database-redis-user"
description: |-
  Creates an user for a redis cluster associated with a public cloud project.
---

# ovh_cloud_project_database_redis_user

Creates an user for a redis cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database" "redis" {
  service_name  = "XXXX"
  engine        = "redis"
  cluster_id    = "ZZZZ"
}

resource "ovh_cloud_project_database_redis_user" "user" {
  service_name  = ovh_cloud_project_database.redis.service_name
  cluster_id    = ovh_cloud_project_database.redis.id
  categories    = ["+@set", "+@sortedset"]
  channels      = ["*"]
  commands	    = ["+get", "-set"]
  keys		      = ["data", "properties"]
  name          = "johndoe"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `categories` - (Optional) Categories of the user.

* `channels` - (Optional: if omit, all channels) Channels of the user.

* `commands` - (Optional) Commands of the user.

* `keys` - (Optional) Keys of the user.

* `name` - (Required, Forces new resource) Name of the user.

## Attributes Reference

The following attributes are exported:

* `categories` - See Argument Reference above.
* `channels` - See Argument Reference above.
* `cluster_id` - See Argument Reference above.
* `commands` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `keys` - See Argument Reference above.
* `name` - See Argument Reference above.
* `password` - (Sensitive) Password of the user.
* `service_name` - See Argument Reference above.
* `status` - Current status of the user.

## Import

OVHcloud Managed redis clusters users can be imported using the `service_name`, `cluster_id` and `id` of the user, separated by "/" E.g.,

```
$ terraform import ovh_cloud_project_database_redis_user.my_user <service_name>/<cluster_id>/<id>