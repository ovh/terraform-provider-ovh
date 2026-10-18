---
subcategory : "Additional IP"
---

# ovh_ip_firewall_rule

Use this resource to manage a rule on an IP firewall.

## Example Usage

```terraform
resource "ovh_ip_firewall_rule" "deny_tcp" {
  ip             = "XXXXXX"
  ip_on_firewall = "XXXXXX"
  sequence       = 0
  action         = "deny"
  protocol       = "tcp"
}
```

With a port range:

```terraform
resource "ovh_ip_firewall_rule" "allow_app_ports" {
  ip             = "XXXXXX"
  ip_on_firewall = "XXXXXX"
  sequence       = 1
  action         = "permit"
  protocol       = "tcp"

  destination_port_range = {
    from = 8000
    to   = 8010
  }
}
```

## Argument Reference

* `ip` - (Required) The IP or the CIDR
* `ip_on_firewall` - (Required) IPv4 address
* `action` - (Required) Possible values for action (deny|permit)
* `protocol` - (Required) Possible values for protocol (ah|esp|gre|icmp|ipv4|tcp|udp)
* `sequence` - (Required) Rule position in the rules array
* `destination_port` - Destination port for your rule. Only with TCP/UDP protocol. Conflicts with `destination_port_range`
* `destination_port_range` - Destination port range for your rule. Only with TCP/UDP protocol. Conflicts with `destination_port`
  * `from` - (Required) Start of the range, between 1 and 65535
  * `to` - (Required) End of the range, between 1 and 65535. Must be greater than `from`: for a single port, use `destination_port`
* `fragments` - Fragments option
* `source` - IPv4 CIDR notation (e.g., 192.0.2.0/24)
* `source_port` - Source port for your rule. Only with TCP/UDP protocol. Conflicts with `source_port_range`
* `source_port_range` - Source port range for your rule. Only with TCP/UDP protocol. Conflicts with `source_port`
  * `from` - (Required) Start of the range, between 1 and 65535
  * `to` - (Required) End of the range, between 1 and 65535. Must be greater than `from`: for a single port, use `source_port`
* `tcp_option` - TCP option on your rule (syn|established)

## Attributes Reference

* `ip` - The IP or the CIDR
* `ip_on_firewall` - IPv4 address
* `state` - Current state of your rule
* `action` - Possible values for action (deny|permit)
* `creation_date` - Creation date of the rule
* `destination` - Destination IP for your rule
* `destination_port` - Destination port for your rule. Only with TCP/UDP protocol
* `destination_port_desc` - String description of field `destination_port` or `destination_port_range` (e.g. `eq 22` or `range 8000 8010`)
* `destination_port_range` - Destination port range for your rule
* `fragments` - Fragments option
* `protocol` - Possible values for protocol (ah|esp|gre|icmp|ipv4|tcp|udp)
* `rule` - Description of the rule
* `sequence` - Rule position in the rules array
* `source` - IPv4 CIDR notation (e.g., 192.0.2.0/24)
* `source_port` - Source port for your rule. Only with TCP/UDP protocol
* `source_port_desc` - String description of field `source_port` or `source_port_range` (e.g. `eq 22` or `range 8000 8010`)
* `source_port_range` - Source port range for your rule
* `tcp_option` - TCP option on your rule (syn|established)

## Import

The resource can be imported using the properties `ip`, `ip_on_firewall` and `sequence`, separated by "|" E.g.,

```bash
$ terraform import ovh_ip_firewall_rule.my_firewall_rule '127.0.0.1|127.0.0.2|0'
```
