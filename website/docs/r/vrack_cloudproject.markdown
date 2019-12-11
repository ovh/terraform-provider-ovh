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
  vrack_id   = "12345"
  project_id = "67890"
}
```

## Argument Reference

The following arguments are supported:

* `vrack_id` - (Required) The id of the vrack. If omitted, the `OVH_VRACK_ID`
    environment variable is used. 
    Note: The use of environment variable is deprecated.

* `project_id` - (Required) The id of the public cloud project. If omitted,
    the `OVH_PROJECT_ID` environment variable is used.
    Note: The use of environment variable is deprecated.
    
## Attributes Reference

The following attributes are exported:

* `vrack_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
