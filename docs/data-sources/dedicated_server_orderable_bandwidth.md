---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_order_bandwidth (Data Source)

Use this data source to get the list of orderable additional bandwidth for a dedicated server associated with your OVHcloud Account.

## Example Usage

```terraform
data "ovh_dedicated_server_orderable_bandwidth" "bp" {
  service_name = "myserver"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your dedicated server.

## Attributes Reference

The following attributes are exported:

* `orderable` - Wether or not additional bandwidth is orderable.
* `platinium` - The list of orderable platinimum bandwidth in mbps.
* `ultimate` - The list of orderable ultimate bandwidth in mbps.
* `premium` - The list of orderable premium bandwidth in mbps.
