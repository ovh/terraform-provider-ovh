---
layout: "ovh"
page_title: "OVH: cloud_project_database_kafka_certificates"
sidebar_current: "docs-ovh-datasource-cloud-project-database-kafka-certificates"
description: |-
  Get information about certificates of a kafka cluster associated with a public cloud project.
---

# ovh_cloud_project_database_kafka_certificates (Data Source)

Use this data source to get information about certificates of a kafka cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_kafka_certificates" "certificates" {
  service_name = "XXX"
  cluster_id   = "YYY"
}

output "certificates_ca" {
  value = data.ovh_cloud_project_database_kafka_certificates.certificates.ca
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

The following attributes are exported:

`id` is set to the md5 sum of the CA. In addition,
the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `ca` - CA certificate used for the service.

