---
subcategory : "Additional IP"
---

# ovh_ip_firewall

Use this resource to manage an IP firewall.

## Example Usage

```hcl
resource "ovh_ip_firewall" "my_firewall" {
  ip             = "XXXXXX"
  ip_on_firewall = "XXXXXX"
}
```

## Argument Reference

* `ip` - (Required) The IP or the CIDR
* `ip_on_firewall` - (Required) IPv4 address
* `enabled ` - Whether firewall should be enabled

## Attributes Reference

* `ip` - The IP or the CIDR
* `ip_on_firewall` - IPv4 address
* `enabled ` - Whether firewall is enabled
* `state` - Current state of your ip on firewall
