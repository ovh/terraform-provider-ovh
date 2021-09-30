---
layout: "ovh"
page_title: "OVH: ovh_iploadbalancing_tcp_route_rule"
sidebar_current: "docs-ovh-resource-iploadbalancing-tcp-route-rule"
description: |-
  Manage rules for TCP route.
---

# ovh_iploadbalancing_tcp_route_rule

Manage rules for TCP route.

## Example Usage

```hcl
resource "ovh_iploadbalancing_tcp_route" "reject" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  weight       = 1
  frontend_id  = 11111

  action {
    type = "reject"
  }
}

resource "ovh_iploadbalancing_tcp_route_rule" "examplerule" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  route_id     = ovh_iploadbalancing_tcp_route.reject.id
  display_name = "Match example.com host"
  field        = "sni"
  match        = "is"
  negate       = false
  pattern      = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `route_id` - (Required) The route to apply this rule
* `display_name` - Human readable name for your rule, this field is for you
* `field` - (Required) Name of the field to match like "protocol" or "host". See "/ipLoadbalancing/{serviceName}/availableRouteRules" for a list of available rules
* `match` - (Required) Matching operator. Not all operators are available for all fields. See "/ipLoadbalancing/{serviceName}/availableRouteRules"
* `negate` - Invert the matching operator effect
* `pattern` - Value to match against this match. Interpretation if this field depends on the match and field
* `sub_field` - Name of sub-field, if applicable. This may be a Cookie or Header name for instance

## Attributes Reference

No additional attribute is exported.
