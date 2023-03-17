---
subcategory : "vRack"
---

# ovh_vrack_iploadbalancing

Attach an IP Load Balancing to a VRack.

## Example Usage

```hcl
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
