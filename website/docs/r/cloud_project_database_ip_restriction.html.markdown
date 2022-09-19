---
layout: "ovh"
page_title: "OVH: cloud_project_database_ip_restriction"
sidebar_current: "docs-ovh-resource-cloud-project-database-ip-restriction"
description: |-
  Creates an IP restriction for a managed database cluster associated with a public cloud project.
---

# ovh_cloud_project_database_ip_restriction

Apply IP restrictions to an OVHcloud Managed Database cluster.

## Example Usage

```hcl
data "ovh_cloud_project_database" "db" {
  service_name = "XXXX"
  engine = "YYYY"
  id  = "ZZZZ"
}

resource "ovh_cloud_project_database_ip_restriction" "iprestriction" {
  service_name = data.ovh_cloud_project_database.db.service_name
  engine       = data.ovh_cloud_project_database.db.engine
  cluster_id   = data.ovh_cloud_project_database.db.id
  ip           = "178.97.6.0/24"
}
```

## Argument Reference

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required, Forces new resource) The engine of the database cluster you want to add an IP restriction. To get a full list of available engine visit.
[public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `ip` - (Required, Forces new resource) Authorized IP.

* `description` - (Optional) Description of the IP restriction.

## Attributes Reference

The following attributes are exported:

* `description` - See Argument Reference above.
* `ip` - See Argument Reference above.
* `status` - Current status of the IP restriction.

## Timeouts

```hcl
resource "ovh_cloud_project_database_ip_restriction" "iprestriction" {
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

OVHcloud Managed database cluster IP restrictions can be imported using the `service_name`, `engine`, `cluster_id` and the `ip`, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_database_ip_restriction.my_ip_restriction service_name/engine/cluster_id/178.97.6.0/24
```
