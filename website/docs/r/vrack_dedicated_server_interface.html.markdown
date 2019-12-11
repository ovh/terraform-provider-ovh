---
layout: "ovh"
page_title: "OVH: vrack_dedicated_server_interface"
sidebar_current: "docs-ovh-resource-vrack-dedicated-server-interface"
description: |-
  Attach a Dedicated Server Network Interface to a VRack.
---

# ovh_vrack_dedicated_server_interface

Attach a Dedicated Server Network Interface to a VRack.

## Example Usage

```hcl
resource "ovh_vrack_dedicated_server_interface" "vdsi" {
  vrack_id   = "12345"
  interface_id = "67890"
}
```

## Argument Reference

The following arguments are supported:

* `vrack_id` - (Required) The id of the vrack.
* `interface_id` - (Required) The id of dedicated server network interface.

## Attributes Reference

The following attributes are exported:

* `vrack_id` - See Argument Reference above.
* `interface_id` - See Argument Reference above.
