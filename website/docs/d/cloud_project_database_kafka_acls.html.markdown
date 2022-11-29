---
layout: "ovh"
page_title: "OVH: cloud_project_database_kafka_acls"
sidebar_current: "docs-ovh-datasource-cloud-project-database-kafka-acls"
description: |-
  Get the list of ACLS of a kafka cluster associated with a public cloud project.
---

# ovh_cloud_project_database_kafka_acls (Data Source)

Use this data source to get the list of ACLs of a kafka cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_kafka_acls" "acls" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "acl_ids" {
  value = data.ovh_cloud_project_database_kafka_acls.acls.acl_ids
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

`id` is set to the md5 sum of the list of all ACL ids. In addition,
the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `acl_ids` - The list of ACLs ids of the kafka cluster associated with the project.
