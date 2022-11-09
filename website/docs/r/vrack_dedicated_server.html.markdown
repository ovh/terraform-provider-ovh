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
  service_name = "XXXX"
  server_id    = "67890"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service name of the vrack. If omitted,
    the `OVH_VRACK_SERVICE` environment variable is used. 

* `server_id` - (Required) The id of the dedicated server. 

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `server_id` - See Argument Reference above.
