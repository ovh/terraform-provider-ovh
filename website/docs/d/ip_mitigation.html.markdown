---
subcategory : "Additional IP"
---

# ovh_ip_mitigation (Data Source)

Use this resource to retrieve information about an IP permanent mitigation.

## Example Usage

```hcl
data "ovh_ip_mitigation" "mitigation_data" {
  ip               = "XXXXXX"
  ip_on_mitigation = "XXXXXX"
}
```

## Argument Reference

* `ip` - (Required) The IP or the CIDR
* `ip_on_mitigation` - (Required) IPv4 address

## Attributes Reference

* `ip` - The IP or the CIDR
* `ip_on_mitigation` - IPv4 address
* `permanent ` - Set on true if the IP is on permanent mitigation
* `state` - Current state of the IP on mitigation
* `auto` - Set on true if the IP is on auto-mitigation
