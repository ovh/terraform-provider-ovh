---
subcategory : "vRack"
---

# ovh_vrack_dedicated_server

Attach a legacy dedicated server to a vRack.

~> **NOTE:** The resource `ovh_vrack_dedicated_server` is intended to be used for legacy dedicated servers.<br />
Dedicated servers that have configurable network interfaces MUST use the resource [`ovh_vrack_dedicated_server_interface`](vrack_dedicated_server_interface.html.markdown) instead.

## Example Usage

```hcl
resource "ovh_vrack_dedicated_server" "vds" {
  service_name = "XXXX"
  server_id    = "67890"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service name of the vrack. If omitted,
    the `OVH_VRACK_SERVICE` environment variable is used.

* `server_id` - (Required) The id of the dedicated server. 

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `server_id` - See Argument Reference above.
