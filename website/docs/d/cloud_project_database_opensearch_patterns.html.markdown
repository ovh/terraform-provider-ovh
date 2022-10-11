---
layout: "ovh"
page_title: "OVH: cloud_project_database_opensearch_patterns"
sidebar_current: "docs-ovh-datasource-cloud-project-database-opensearch-patterns"
description: |-
  Get the list of patterns of a opensearch cluster associated with a public cloud project.
---

# ovh_cloud_project_database_opensearch_patterns (Data Source)

Use this data source to get the list of pattern of a opensearch cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database_opensearch_patterns" "patterns" {
  service_name  = "XXX"
  cluster_id    = "YYY"
}

output "pattern_ids" {
  value = data.ovh_cloud_project_database_opensearch_patterns.patterns.pattern_ids
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

`id` is set to the md5 sum of the list of all patterns ids. In addition,
the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `pattern_ids` - The list of patterns ids of the opensearch cluster associated with the project.
