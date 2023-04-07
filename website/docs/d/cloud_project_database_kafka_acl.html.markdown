---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_kafka_acl (Data Source)

Use this data source to get information about an ACL of a kafka cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_kafka_acl" "acl" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  id            = "ZZZ"
}

output "acl_permission" {
  value = data.ovh_cloud_project_database_kafka_acl.acl.permission
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `id` - (Required) ACL ID

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `id` - See Argument Reference above.
* `permission` - Permission to give to this username on this topic.
* `service_name` - See Argument Reference above.
* `topic` - Topic affected by this ACL.
* `username` - Username affected by this ACL.
