---
layout: "ovh"
page_title: "OVH: cloud_project_database_kafka_user_acces"
sidebar_current: "docs-ovh-datasource-cloud-project-database-kafka-user-acces"
description: |-
  Get information about user acces of a kafka cluster associated with a public cloud project.
---

# ovh_cloud_project_database_kafka_user_access (Data Source)

Use this data source to get information about user acces of a kafka cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_kafka_user_access" "access" {
  service_name = "XXX"
  cluster_id   = "YYY"
  user_id      = "ZZZ"
}

output "access_cert" {
  value = data.ovh_cloud_project_database_kafka_user_access.access.cert
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `user_id` - (Required) User ID

## Attributes Reference

`id` is set to the md5 sum of the Cert. In addition,
the following attributes are exported:

* `cert` - User cert.
* `cluster_id` - See Argument Reference above.
* `key` - (Sensitive) User key for the cert.
* `service_name` - See Argument Reference above.
* `user_id` - See Argument Reference above.
