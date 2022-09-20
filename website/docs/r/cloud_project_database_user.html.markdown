---
layout: "ovh"
page_title: "OVH: cloud_project_database_user"
sidebar_current: "docs-ovh-resource-cloud-project-database-user"
description: |-
  Creates an user for a database cluster associated with a public cloud project.
---

# ovh_cloud_project_database_user

Creates an user for a database cluster associated with a public cloud project.

With this resource you can create a user for the following database engine:

  * `cassandra`
  * `kafka`
  * `kafkaConnect`
  * `mysql`

## Example Usage

```hcl
data "ovh_cloud_project_database" "db" {
  service_name  = "XXXX"
  engine        = "YYYY"
  id            = "ZZZZ"
}

resource "ovh_cloud_project_database_user" "user" {
  service_name  = data.ovh_cloud_project_database.db.service_name
  engine        = data.ovh_cloud_project_database.db.engine
  cluster_id    = data.ovh_cloud_project_database.db.id
  name          = "johndoe"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required, Forces new resource) The engine of the database cluster you want to add. To get a full list of available engine visit :
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).\
Available engines for this resource (other have specific resource):
  * `cassandra`
  * `kafka`
  * `kafkaConnect`
  * `mysql`

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `name` - (Required, Forces new resource) Name of the user.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `password` - (Sensitive) Password of the user.
* `service_name` - See Argument Reference above.
* `status` - Current status of the user.
* `name` - See Argument Reference above.

## Timeouts

```hcl
resource "ovh_cloud_project_database_user" "user" {
  # ...

  timeouts {
    create = "1h"
    delete = "45m"
  }
}
```
* `create` - (Default 20m)
* `delete` - (Default 20m)

## Import

OVHcloud Managed database clusters users can be imported using the `service_name`, `engine`, `cluster_id` and `id` of the user, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_database_user.my_user service_name/engine/cluster_id/id
```
