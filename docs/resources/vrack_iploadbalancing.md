---
subcategory : "vRack"
---

# ovh_vrack_iploadbalancing

Attach an IP Load Balancing to a VRack.

## Example Usage

```terraform
resource "ovh_vrack_iploadbalancing" "viplb" {
  service_name     = "xxx"
  ip_loadbalancing = "yyy"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the vrack.
* `ip_loadbalancing` - (Required) The id of the IP Load Balancing.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `ip_loadbalancing` - See Argument Reference above.

## Import

A vRack IP Load Balancing attachment can be imported using the `service_name` and `ip_loadbalancing`, separated by "/" E.g.,

```bash
$ terraform import ovh_vrack_iploadbalancing.viplb service_name/ip_loadbalancing
```
