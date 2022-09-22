---
layout: "ovh"
page_title: "OVH: cloud_project_database_m3db_namespaces"
sidebar_current: "docs-ovh-datasource-cloud-project-database-m3db-namespaces"
description: |-
  Get the list of namespaces of a m3db cluster associated with a public cloud project.
---

# ovh_cloud_project_database_m3db_namespaces (Data Source)

Use this data source to get the list of namespaces of a m3db cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_m3db_namespaces" "namespaces" {
  service_name = "XXX"
  cluster_id   = "YYY"
}

output "namespace_ids" {
  value = data.ovh_cloud_project_database_m3db_namespaces.namespaces.namespace_ids
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

`id` is set to the md5 sum of the list of all namespaces ids. In addition,
the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `namespace_ids` - The list of namespaces ids of the m3db cluster associated with the project.
