---
layout: "ovh"
page_title: "OVH: ovh_iploadbalancing_vrack_networks"
sidebar_current: "docs-ovh-iploadbalancing-vrack-networks"
description: |-
  Get the list of Vrack network ids available for your IPLoadbalancer associated with your OVH account.
---

# ovh_iploadbalancing_vrack_networks

Use this data source to get the list of Vrack network ids available for your IPLoadbalancer associated with your OVH account.

## Example Usage

```hcl
data ovh_iploadbalancing_vrack_networks "lb_networks" {
  service_name = "xxx"
  subnet       = "10.0.0.0/24"
}
```

## Argument Reference


* `service_name` - (Required) The internal name of your IP load balancing

* `subnet` - Filters networks on the subnet.

* `vlan_id` - Filters networks on the vlan id.


## Attributes Reference

The following attributes are exported:

* `result` - The list of vrack network ids.

