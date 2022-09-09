---
layout: "ovh"
page_title: "OVH: cloud_project_database_opensearch_pattern"
sidebar_current: "docs-ovh-resource-cloud-project-database-opensearch-pattern"
description: |-
  Creates a pattern for a opensearch cluster associated with a public cloud project.
---

# ovh_cloud_project_database_opensearch_pattern

Creates a pattern for a opensearch cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database" "opensearch" {
  service_name  = "XXX"
  engine        = "opensearch"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_opensearch_pattern" "pattern" {
  service_name = data.ovh_cloud_project_database.opensearch.service_name
  cluster_id   = data.ovh_cloud_project_database.opensearch.id
  max_index_count = 2
  pattern = "logs_*"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `max_index_count` - (Optional, Forces new resource) Maximum number of index for this pattern.

* `pattern` - (Required, Forces new resource) Pattern format.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `id` - ID of the pattern.
* `max_index_count` - See Argument Reference above.
* `pattern` - See Argument Reference above.
* `service_name` - See Argument Reference above.

## Timeouts

```hcl
resource "ovh_cloud_project_database_opensearch_pattern" "pattern" {
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

OVHcloud Managed opensearch clusters patterns can be imported using the `service_name`, `cluster_id` and `id` of the pattern, separated by "/" E.g.,

```
$ terraform import ovh_cloud_project_database_opensearch_pattern.my_pattern service_name/cluster_id/id