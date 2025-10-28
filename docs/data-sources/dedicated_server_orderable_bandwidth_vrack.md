---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_orderable_bandwidth_vrack (Data Source)

Use this data source to get the orderable vrack bandwidth information about a dedicated server associated with your OVHcloud Account.

## Example Usage

```terraform
data "ovh_dedicated_server_orderable_bandwidth_vrack" "spec" {
  service_name = "myserver"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your dedicated server.

## Attributes Reference

The following attributes are exported:

* `orderable` - Whether or not additional bandwidth is orderable
* `vrack` - Additional orderable vrack bandwidth