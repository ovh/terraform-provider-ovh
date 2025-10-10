---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_ip_restriction

~> **DEPRECATED** Use `ip_restrictions` field in `ovh_cloud_project_database` resource instead. Using the `ovh_cloud_project_database_ip_restriction` resource to add an IP restriction to a `ovh_cloud_project_database` resource will cause its update on each apply

~> **NOTE** This resource is only compatible with the following engines: **cassandra, grafana, kafka, kakfaConnect, m3db, mongodb, mysql, opensearch, postgres, redis**

Apply IP restrictions to an OVHcloud Managed Database cluster.

## Example Usage

```terraform
data "ovh_cloud_project_database" "db" {
  service_name = "XXXX"
  engine       = "YYYY"
  id           = "ZZZZ"
}

resource "ovh_cloud_project_database_ip_restriction" "ip_restriction" {
  service_name = data.ovh_cloud_project_database.db.service_name
  engine       = data.ovh_cloud_project_database.db.engine
  cluster_id   = data.ovh_cloud_project_database.db.id
  ip           = "178.97.6.0/24"
}
```

## Argument Reference

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required, Forces new resource) The engine of the database cluster you want to add an IP restriction. To get a full list of available engine visit. [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases).

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `ip` - (Required, Forces new resource) Authorized IP.

* `description` - (Optional) Description of the IP restriction.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `description` - See Argument Reference above.
* `engine` - See Argument Reference above.
* `ip` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `status` - Current status of the IP restriction.

## Timeouts

```terraform
resource "ovh_cloud_project_database_ip_restriction" "ip_restriction" {
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
