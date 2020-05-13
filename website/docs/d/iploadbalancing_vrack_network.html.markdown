---
layout: "ovh"
page_title: "OVH: ovh_iploadbalancing_vrack_network"
sidebar_current: "docs-ovh-iploadbalancing-vrack-network"
description: |-
  Get the details of Vrack network available for your IPLoadbalancer associated with your OVH account.
---

# ovh_iploadbalancing_vrack_network

Use this data source to get the details of Vrack network available for your IPLoadbalancer associated with your OVH account.

## Example Usage

```hcl
data ovh_iploadbalancing_vrack_network "lb_network" {
  service_name = "xxx"
  vrack_network_id = "yyy"
}
```

## Argument Reference


* `service_name` - (Required) The internal name of your IP load balancing

* `vrack_network_id` - (Required) Internal Load Balancer identifier of the vRack private network


## Attributes Reference

The following attributes are exported:

* `display_name` - Human readable name for your vrack network
* `nat_ip` - An IP block used as a pool of IPs by this Load Balancer to connect to the servers in this private network. The blck must be in the private network and reserved for the Load Balancer
* `subnet` - IP block of the private network in the vRack
* `vlan` - VLAN of the private network in the vRack. 0 if the private network is not in a VLAN

