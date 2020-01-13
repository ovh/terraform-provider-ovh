---
layout: "ovh"
page_title: "OVH: vrack_iploadbalancing"
sidebar_current: "docs-ovh-resource-vrack-ip-loadbalancing"
description: |-
  Attach a IP Loadbalanging to a VRack.
---

# ovh_vrack_iploadbalancing

Attach a ip loadbalancing to a VRack.

## Example Usage

```hcl
resource "ovh_vrack_iploadbalancing" "viplb" {
  service_name   = "xxx"
  ip_loadbalancing = "yyy"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the vrack.
* `ip_loadbalancing` - (Required) The id of the ip loadbalancing. 

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `ip_loadbalancing` - See Argument Reference above.
