---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing_http_route_rule

Manage rules for HTTP route.

## Example Usage

Route which redirect all URL to HTTPs for example.com (Vhost).

```hcl
resource "ovh_iploadbalancing_http_route" "httpsredirect" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  display_name = "Redirect to HTTPS"
  weight       = 1
  frontend_id  = 11111

  action {
    status = 302
    target = "https://$${host}$${path}$${arguments}"
    type   = "redirect"
  }
}

resource "ovh_iploadbalancing_http_route_rule" "examplerule" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  route_id     = "${ovh_iploadbalancing_http_route.httpsredirect.id}"
  display_name = "Match example.com host"
  field        = "host"
  match        = "is"
  negate       = false
  pattern      = "example.com"
}
```

Rule which match a specific header (same effect as the host match above).

```hcl
resource "ovh_iploadbalancing_http_route_rule" "examplerule" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  route_id     = "${ovh_iploadbalancing_http_route.httpsredirect.id}"
  display_name = "Match example.com Host header"
  field        = "headers"
  match        = "is"
  negate       = false
  pattern      = "example.com"
  sub_field    = "Host"
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

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `route_id` - See Argument Reference above.
* `display_name` - See Argument Reference above.
* `field` - See Argument Reference above.
* `match` - See Argument Reference above.
* `negate` - See Argument Reference above.
* `pattern` - See Argument Reference above.
* `sub_field` - See Argument Reference above.


## Import 

HTTP route rule can be imported using the following format `service_name`, the `id` of the route and the `id` of the rule separated by "/" e.g.

```bash
$ terraform import ovh_iploadbalancing_http_route_rule.examplerule service_name/route_id/rule_id
```