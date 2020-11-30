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

* `vrack_id` - (Optional) Deprecated. The id of the vrack. If omitted,
    the `OVH_VRACK_ID` environment variable is used.
    One of `service_name` or `vrack_id` is required. Conflits with `service_name`.

* `service_name` - (Optional) The id of the vrack. If omitted,
    the `OVH_VRACK_SERVICE` environment variable is used. 
    One of `service_name` or `vrack_id` is required. Conflits with `vrack_id`.

* `project_id` - (Required) The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

    
## Attributes Reference

The following attributes are exported:

* `vrack_id` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `project_id` - See Argument Reference above.
