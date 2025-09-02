---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_order_bandwidth_vrack (Data Source)

Use this data source to get the list of orderable additional vrack bandwidth for a dedicated server associated with your OVHcloud Account.

## Example Usage

```terraform
data "ovh_dedicated_server_orderable_bandwidth_vrack" "bp" {
  service_name = "myserver"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your dedicated server.

## Attributes Reference

The following attributes are exported:

* `orderable` - Wether or not additional bandwidth is orderable.
* `vrack` - The list of orderable vrack bandwidth in mbps.
