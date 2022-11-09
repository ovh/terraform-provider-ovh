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
data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

resource "ovh_vrack_dedicated_server_interface" "vdsi" {
  service_name = "pn-xxxxxxx" #name of the vrack
  interface_id = data.ovh_dedicated_server.server.enabled_vrack_vnis[0]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the vrack. If omitted,
    the `OVH_VRACK_SERVICE` environment variable is used. 

* `interface_id` - (Required) The id of dedicated server network interface.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `interface_id` - See Argument Reference above.
