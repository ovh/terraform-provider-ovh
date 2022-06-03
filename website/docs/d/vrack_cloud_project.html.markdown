---
layout: "ovh"
page_title: "OVH: ovh_vrack_cloud_project"
sidebar_current: "docs-ovh-vrack-cloud-project"
description: |-
  Get the list of cloud projects attach to your Vrack ids.
---

# ovh_vrack_cloud_project (Data Source)

Use this data source to get the list of cloud projects attach to your Vrack IDs.

## Example Usage

```hcl
data "ovh_vrack_cloud_project" "vrack" {
   service_name = "XXXXXX"
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your Vrack.

## Attributes Reference

The following attributes are exported:

*  - The list of cloud projects attach to your Vrack ids.
