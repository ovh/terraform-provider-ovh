---
layout: "ovh"
page_title: "OVH: cloud_regions"
sidebar_current: "docs-ovh-datasource-cloud-regions"
description: |-
  Get the list of regions associated with a public cloud project.
---

# ovh_cloud_regions

Use this data source to get the regions of a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_regions" "regions" {
  project_id = "XXXXXX"
}
```

## Argument Reference


* `project_id` - (Required) The id of the public cloud project. If omitted,
    the `OVH_PROJECT_ID` environment variable is used.


## Attributes Reference

`id` is set to the ID of the project. In addition, the following attributes
are exported:

* `names` - The list of regions associated with the project
