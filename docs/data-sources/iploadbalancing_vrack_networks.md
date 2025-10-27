---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing_vrack_networks (Data Source)

Use this data source to get the list of Vrack network ids available for your IPLoadbalancer associated with your OVHcloud account.

## Example Usage

```terraform
data "ovh_iploadbalancing_vrack_networks" "lb_networks" {
  service_name = "XXXXXX"
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
