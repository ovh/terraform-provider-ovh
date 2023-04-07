---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing_vrack_network

Manage a vrack network for your IP Loadbalancing service.

## Example Usage

```hcl
data ovh_iploadbalancing "iplb" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
}

resource "ovh_vrack_iploadbalancing" "viplb" {
  service_name     = "xxx"
  ip_loadbalancing = data.ovh_iploadbalancing.iplb.service_name
}

resource ovh_iploadbalancing_vrack_network "network" {
  service_name = ovh_vrack_iploadbalancing.viplb.ip_loadbalancing
  subnet       = "10.0.0.0/16"
  vlan         = 1
  nat_ip       = "10.0.0.0/27"
  display_name = "mynetwork"
}

resource "ovh_iploadbalancing_tcp_farm" "testfarm" {
  service_name     = ovh_iploadbalancing_vrack_network.network.service_name
  display_name     = "mytcpbackends"
  port             = 80
  vrack_network_id = ovh_iploadbalancing_vrack_network.network.vrack_network_id
  zone             = tolist(data.ovh_iploadbalancing.iplb.zone)[0]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `display_name` - Human readable name for your vrack network
* `farm_id` - This attribute is there for documentation purpose only and isnt passed to the OVHcloud API as it may conflicts with http/tcp farms `vrack_network_id` attribute
* `nat_ip` - (Required) An IP block used as a pool of IPs by this Load Balancer to connect to the servers in this private network. The blck must be in the private network and reserved for the Load Balancer
* `subnet` - (Required) IP block of the private network in the vRack
* `vlan` - VLAN of the private network in the vRack. 0 if the private network is not in a VLAN

## Attributes Reference

The following attributes are exported:

* `vrack_network_id` - (Required) Internal Load Balancer identifier of the vRack private network
