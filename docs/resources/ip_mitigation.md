---
subcategory : "Additional IP"
---

# ovh_ip_mitigation

Use this resource to manage an IP permanent mitigation.

## Example Usage

```terraform
resource "ovh_ip_mitigation" "mitigation" {
  ip               = "XXXXXX"
  ip_on_mitigation = "XXXXXX"
}
```

## Argument Reference

* `ip` - (Required) The IP or the CIDR
* `ip_on_mitigation` - (Required) IPv4 address
* `permanent ` - Deprecated, has no effect

## Attributes Reference

* `ip` - The IP or the CIDR
* `ip_on_mitigation` - IPv4 address
* `permanent ` - (Deprecated) Set on true if the IP is on permanent mitigation
* `state` - Current state of the IP on mitigation
* `auto` - Set on true if the IP is on auto-mitigation
