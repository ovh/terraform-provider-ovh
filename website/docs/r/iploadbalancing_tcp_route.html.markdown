---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing_tcp_route

Manage TCP route for a loadbalancer service

## Example Usage

```hcl
resource "ovh_iploadbalancing_tcp_route" "tcpreject" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  weight = 1

  action {
    type = "reject"
  }
}
```

## Argument Reference

The following arguments are supported:

* `action` - (Required) Action triggered when all rules match
   * `target` - Farm ID for "farm" action type, empty for others.
   * `type` - (Required) Action to trigger if all the rules of this route matches
* `display_name` - Human readable name for your route, this field is for you
* `frontend_id` - Route traffic for this frontend
* `service_name` - (Required) The internal name of your IP load balancing
* `weight` - Route priority ([0..255]). 0 if null. Highest priority routes are evaluated first. Only the first matching route will trigger an action

## Attributes Reference

In addition, the following attributes are exported:

* `status` - Route status. Routes in "ok" state are ready to operate
* `rules` - List of rules to match to trigger action
   * `field` - Name of the field to match like "protocol" or "host" "/ipLoadbalancing/{serviceName}/route/availableRules" for a list of available rules
   * `match` - Matching operator. Not all operators are available for all fields. See "availableRules"
   * `negate`- Invert the matching operator effect
   * `pattern` - Value to match against this match. Interpretation if this field depends on the match and field
   * `rule_id` - Id of your rule
   * `sub_field` - Name of sub-field, if applicable. This may be a Cookie or Header name for instance

## Import 

TCP route can be imported using the following format `service_name` and the `id` of the route separated by "/" e.g.

```bash
$ terraform import ovh_iploadbalancing_tcp_route.tcpreject service_name/route_id
```