---
subcategory : "vRack"
---

# ovh_vrack_dedicated_server_interface

Attach a Dedicated Server Network Interface to a vRack.

~> **NOTE:** The resource `ovh_vrack_dedicated_server_interface` is intended to be used for dedicated servers that have configurable network interfaces.<br />
Legacy Dedicated servers that do not have configurable network interfaces MUST use the resource [`ovh_vrack_dedicated_server`](vrack_dedicated_server.html.markdown) instead.

## Example Usage

```hcl
data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

resource "ovh_vrack_dedicated_server_interface" "vdsi" {
  service_name = "pn-xxxxxxx" #name of the vrack
  interface_id = data.ovh_dedicated_server.server.enabled_vrack_vnis[0]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the vrack. If omitted,
    the `OVH_VRACK_SERVICE` environment variable is used.

* `interface_id` - (Required) The id of dedicated server network interface.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `interface_id` - See Argument Reference above.
