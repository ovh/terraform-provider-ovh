---
subcategory : "Additional IP"
---

# ovh_ip_firewall_rule (Data Source)

Use this data source to retrieve information about a rule on an IP firewall.

## Example Usage

```terraform
data "ovh_ip_firewall_rule" "my_firewall_rule" {
  ip             = "XXXXXX"
  ip_on_firewall = "XXXXXX"
  sequence       = 0
}
```

## Argument Reference

* `ip` - (Required) The IP or the CIDR
* `ip_on_firewall` - (Required) IPv4 address
* `sequence` - (Required) Rule position in the rules array

## Attributes Reference

* `ip` - The IP or the CIDR
* `ip_on_firewall` - IPv4 address
* `state` - Current state of your rule
* `action` - Possible values for action (deny|permit)
* `creation_date` - Creation date of the rule
* `destination` - Destination IP for your rule
* `destination_port` - Destination port or port range for your rule, as returned by the API (e.g. `eq 22` or `range 8000 8010`). Only with TCP/UDP protocol
* `destination_port_range` - Destination port range for your rule, when it has one
  * `from` - Start of the range
  * `to` - End of the range
* `fragments` - Fragments option
* `protocol` - Possible values for protocol (ah|esp|gre|icmp|ipv4|tcp|udp)
* `rule` - Description of the rule
* `sequence` - Rule position in the rules array
* `source` - IPv4 CIDR notation (e.g., 192.0.2.0/24)
* `source_port` - Source port or port range for your rule, as returned by the API (e.g. `eq 22` or `range 8000 8010`). Only with TCP/UDP protocol
* `source_port_range` - Source port range for your rule, when it has one
  * `from` - Start of the range
  * `to` - End of the range
* `tcp_option` - TCP option on your rule (syn|established)
