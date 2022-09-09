---
layout: "ovh"
page_title: "OVH: cloud_project_database_mongodb_user"
sidebar_current: "docs-ovh-resource-cloud-project-database-mongodb-user"
description: |-
  Creates an user for a mongodb cluster associated with a public cloud project.
---

# ovh_cloud_project_database_mongodb_user

Creates an user for a mongodb cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database" "mongodb" {
  service_name  = "XXX"
  engine        = "mongodb"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_mongodb_user" "user" {
  service_name  = data.ovh_cloud_project_database.mongodb.service_name
  cluster_id    = data.ovh_cloud_project_database.mongodb.id
  name          = "johndoe"
  roles         = ["backup", "readAnyDatabase"]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `name` - (Required, Forces new resource) Name of the user.

* `roles` - (Optional: if omit, default role) Roles the user belongs to. Possible values:
  * `backup`
  * `dbAdminAnyDatabase`
  * `readAnyDatabase`
  * `readWriteAnyDatabase`
  * `restore`
  * `userAdminAnyDatabase`

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `name` - See Argument Reference above.
* `password` - (Sensitive) Password of the user.
* `roles` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `status` - Current status of the user.
* `name` - Name of the user with the authentication database in the format name@authDB

## Timeouts

```hcl
resource "ovh_cloud_project_database_mongodb_user" "user" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}
```
* `create` - (Default 20m)
* `update` - (Default 20m)
* `delete` - (Default 20m)

## Import

OVHcloud Managed mongodb clusters users can be imported using the `service_name`, `cluster_id` and `id` of the user, separated by "/" E.g.,

```
$ terraform import ovh_cloud_project_database_mongodb_user.my_user service_name/cluster_id/id