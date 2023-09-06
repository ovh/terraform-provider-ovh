---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing_failover_ip_attach (Data Source)

Use this data source to check failover IP routed to this IPLB.

## Example Usage

```hcl
data "ovh_iploadbalancing_failover_ip_attach" "failoverip" {
 service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
 ip = "x.x.x.x/y"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `ip` - (Required) IP to move

## Attributes Reference

No attributes exported.
