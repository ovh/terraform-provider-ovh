---
subcategory : "Additional IP"
---

# ovh_ip_firewall

Use this resource to manage a rule on an IP firewall.

## Example Usage

```hcl
resource "ovh_ip_firewall_rule" "myfirewallrule" {
  ip             = "XXXXXX"
  ip_on_firewall = "XXXXXX"
  sequence       = 0
  action         = "deny"
  protocol       = "tcp"
}
```

## Argument Reference

* `ip` - (Required) The IP or the CIDR
* `ip_on_firewall` - (Required) IPv4 address
* `action` - (Required) Possible values for action (deny|permit)
* `protocol` - (Required) Possible values for protocol (ah|esp|gre|icmp|ipv4|tcp|udp)
* `sequence` - (Required) Rule position in the rules array
* `destination_port` - Destination port for your rule. Only with TCP/UDP protocol
* `fragments` - Fragments option
* `source` - IPv4 CIDR notation (e.g., 192.0.2.0/24)
* `tcp_option` - TCP option on your rule (syn|established)

## Attributes Reference

* `ip` - The IP or the CIDR
* `ip_on_firewall` - IPv4 address
* `state` - Current state of your rule
* `action` - Possible values for action (deny|permit)
* `creation_date` - Creation date of the rule
* `destination` - Destination IP for your rule
* `destination_port` - Destination port for your rule. Only with TCP/UDP protocol
* `fragments` - Fragments option
* `protocol` - Possible values for protocol (ah|esp|gre|icmp|ipv4|tcp|udp)
* `rule` - Description of the rule
* `sequence` - Rule position in the rules array
* `source` - IPv4 CIDR notation (e.g., 192.0.2.0/24)
* `source_port` - Source port for your rule. Only with TCP/UDP protocol
* `tcp_option` - TCP option on your rule (syn|established)