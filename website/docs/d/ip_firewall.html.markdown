---
subcategory : "Additional IP"
---

# ovh_ip_firewall (Data Source)

Use this data source to retrieve information about an IP firewall.

## Example Usage

```hcl
data "ovh_ip_firewall" "myfirewall" {
  ip             = "XXXXXX"
  ip_on_firewall = "XXXXXX"
}
```

## Argument Reference

* `ip` - (Required) The IP or the CIDR
* `ip_on_firewall` - (Required) IPv4 address

## Attributes Reference

* `ip` - The IP or the CIDR
* `ip_on_firewall` - IPv4 address
* `enabled ` - Whether firewall is enabled
* `state` - Current state of your ip on firewall