---
layout: "ovh"
page_title: "OVH: cloud_project_database_opensearch_pattern"
sidebar_current: "docs-ovh-datasource-cloud-project-database-opensearch-pattern"
description: |-
  Get information about a pattern of a opensearch cluster associated with a public cloud project.
---

# ovh_cloud_project_database_opensearch_pattern (Data Source)

Use this data source to get information about a pattern of a opensearch cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_opensearch_pattern" "pattern" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  id            = "ZZZ"
}

output "pattern_pattern" {
  value = data.ovh_cloud_project_database_opensearch_pattern.pattern.pattern
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `id` - (Required) Pattern ID.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `id` - See Argument Reference above.
* `max_index_count` - Maximum number of index for this pattern.
* `pattern` - Pattern format.
* `service_name` - Current status of the pattern.
