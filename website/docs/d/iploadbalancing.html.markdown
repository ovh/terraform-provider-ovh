---
layout: "ovh"
page_title: "OVH: iploadbalancing"
sidebar_current: "docs-ovh-datasource-iploadbalancing"
description: |-
  Get information & status of an IP Load Balancing product.
---

# ovh_iploadbalancing

Use this data source to retrieve information about an IP Load Balancing product

## Example Usage

```hcl
data "ovh_iploadbalancing" "lb" {
   service_name = "xxx"
   state        = "ok"
}
```

## Argument Reference

* `ipv6` - The IPV6 associated to your IP load balancing

* `ipv4` - The IPV4 associated to your IP load balancing

* `zone` - Location where your service is. This takes an array of values.

* `offer` - The offer of your IP load balancing

* `service_name` - The internal name of your IP load balancing

* `ip_loadbalancing` - Your IP load balancing

* `state` - Current state of your IP. Can take any of the following value:
"blacklisted", "deleted", "free", "ok", "quarantined", "suspended"

* `vrack_eligibility` - Vrack eligibility. Takes a boolean value.

* `vrack_name` - Name of the vRack on which the current Load Balancer is
attached to, as it is named on vRack product

* `ssl_configuration` - Modern oldest compatible clients : Firefox 27, Chrome 30,
IE 11 on Windows 7, Edge, Opera 17, Safari 9, Android 5.0, and Java 8.
Intermediate oldest compatible clients : Firefox 1, Chrome 1, IE 7, Opera 5,
Safari 1, Windows XP IE8, Android 2.3, Java 7.
Can take any of the following value: "intermediate", "modern"

* `display_name` - the name displayed in ManagerV6 for your iplb (max 50 chars)

## Attributes Reference

`id` is set to the service_name of your IP load balancing
In addition, the following attributes are exported:

* `metrics_token` - The metrics token associated with your IP load balancing
This attribute is sensitive.

* `orderable_zone` - Available additional zone for your Load Balancer
  * `name` - The zone three letter code
  * `plan_code` - The billing planCode for this zone
