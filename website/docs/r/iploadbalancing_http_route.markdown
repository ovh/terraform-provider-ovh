---
layout: "ovh"
page_title: "OVH: ovh_iploadbalancing_http_route"
sidebar_current: "docs-ovh-resource-iploadbalancing-http-route-x"
description: |-
  Manage http route for a loadbalancer service.
---

# ovh_iploadbalancing_http_route

Manage http route for a loadbalancer service

## Example Usage

Route which redirect all url to https.

```hcl
resource "ovh_iploadbalancing_http_route" "httpsredirect" {
  service_name = "loadbalancer-xxxxxxxxxxxxxxxxxx"
  display_name = "Redirect to HTTPS"
  weight = 1

  action {
    status = 302
    target = "https://${host}${path}${arguments}"
    type = "redirect"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `display_name` - Human readable name for your route, this field is for you
* `weight` - Route priority ([0..255]). 0 if null. Highest priority routes are evaluated first. Only the first matching route will trigger an action
* `action.status` - HTTP status code for "redirect" and "reject" actions
* `action.target` - Farm ID for "farm" action type or URL template for "redirect" action. You may use ${uri}, ${protocol}, ${host}, ${port} and ${path} variables in redirect target
* `action.type` - (Required) Action to trigger if all the rules of this route matches
* `frontend_id` - Route traffic for this frontend

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `display_name` - See Argument Reference above.
* `weight` - See Argument Reference above.
* `action.status` - See Argument Reference above.
* `action.target` - See Argument Reference above.
* `action.type` - See Argument Reference above.
* `frontend_id` - See Argument Reference above.
