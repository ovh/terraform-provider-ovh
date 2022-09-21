---
layout: "ovh"
page_title: "OVH: cloud_project_regions"
sidebar_current: "docs-ovh-datasource-cloud-project-regions"
description: |-
  Get the list of regions associated with a public cloud project.
---

# ovh_cloud_project_regions (Data Source)

Use this data source to get the regions of a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_regions" "regions" {
  service_name    = "XXXXXX"
  has_services_up = ["network"]
}
```

## Argument Reference


* `service_name` - (Required) The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. 

* `has_services_up` - (Optional) List of services which has to be UP in regions.
   Example: "image", "instance", "network", "storage", "volume", "workflow", ...
   If left blank, returns all regions associated with the service_name.


## Attributes Reference

`id` is set to the ID of the project. In addition, the following attributes
are exported:

* `names` - The list of regions associated with the project, filtered by services UP.
