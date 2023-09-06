---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing_failover_ip_attach

Attach an IP for a loadbalancer service.

## Example Usage

```hcl
resource "ovh_iploadbalancing_failover_ip_attach" "failoverip" {
 service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
 ip = "x.x.x.x/y"
 to = "loadbalancer-xxxxxxxxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `ip` - (Required) IP to move
* `to` - (Required) Service destination
* `nexthop` - Nexthop of destination service

## Attributes Reference

No attributes exported.
