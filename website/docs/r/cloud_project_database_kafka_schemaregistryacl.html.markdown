---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_kafka_schemaregistryacl

Creates a schema registry ACL for a Kafka cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database" "kafka" {
  service_name  = "XXX"
  engine        = "kafka"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_kafka_schemaregistryacl" "schemaRegistryAcl" {
  service_name    = data.ovh_cloud_project_database.kafka.service_name
  cluster_id      = data.ovh_cloud_project_database.kafka.id
  permission      = "schema_registry_read"
  resource        = "Subject:myResource"
  username        = "johndoe"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `permission` - (Required, Forces new resource) Permission to give to this username on this resource.
Available permissions:
  * `schema_registry_read`
  * `schema_registry_write`

* `resource` - (Required, Forces new resource) Resource affected by this schema registry ACL.

* `username` - (Required, Forces new resource) Username affected by this schema registry ACL.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `id` - ID of the ACL.
* `permission` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `resource` - See Argument Reference above.
* `username` - See Argument Reference above.

## Timeouts

```hcl
resource "ovh_cloud_project_database_kafka_schemaregistryacl" "schemaRegistryAcl" {
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

OVHcloud Managed Kafka clusters schema registry ACLs can be imported using the `service_name`, `cluster_id` and `id` of the schema registry ACL, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_database_kafka_schemaregistryacl.my_schemaRegistryAcl service_name/cluster_id/id
```
