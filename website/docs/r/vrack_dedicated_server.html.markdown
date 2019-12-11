---
layout: "ovh"
page_title: "OVH: vrack_dedicated_server"
sidebar_current: "docs-ovh-resource-vrack-dedicated-server"
description: |-
  Attach a Dedicated Server to a VRack.
---

# ovh_vrack_dedicated_server

Attach a dedicated server to a VRack.

## Example Usage

```hcl
resource "ovh_vrack_dedicated_server" "vds" {
  vrack_id   = "12345"
  server_id = "67890"
}
```

## Argument Reference

The following arguments are supported:

* `vrack_id` - (Required) The id of the vrack.
* `server_id` - (Required) The id of the dedicated server. 

## Attributes Reference

The following attributes are exported:

* `vrack_id` - See Argument Reference above.
* `server_id` - See Argument Reference above.
