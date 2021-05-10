---
layout: "ovh"
page_title: "OVH: vrack_cloudproject"
sidebar_current: "docs-ovh-resource-vrack-cloudproject"
description: |-
  Attach a Public Cloud Project to a VRack.
---

# ovh_vrack_cloudproject

Attach a Public Cloud Project to a VRack.

## Example Usage

```hcl
resource "ovh_vrack_cloudproject" "vcp" {
  service_name = "12345"
  project_id   = "67890"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the vrack. If omitted,
    the `OVH_VRACK_SERVICE` environment variable is used. 

* `project_id` - (Required) The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

    
## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `project_id` - See Argument Reference above.
