---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_orderable_bandwidth (Data Source)

Use this data source to get the list of orderable additional bandwidth for a dedicated server associated with your OVHcloud Account.

## Example Usage

```terraform
data "ovh_dedicated_server_orderable_bandwidth" "spec" {
  service_name = "myserver"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your dedicated server.

## Attributes Reference

The following attributes are exported:

* `orderable` - Whether or not additional bandwidth is orderable
* `platinium` - Additional orderable platinium bandwidth
* `ultimate` - Additional orderable ultimate bandwidth
* `premium` - Additional orderable premium bandwidth