---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing_udp_frontend

Creates a backend server group (frontend) to be used by loadbalancing frontend(s)

## Example Usage

```hcl
data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_udp_frontend" "testfrontend" {
  service_name = data.ovh_iploadbalancing.lb.service_name
  display_name = "ingress-8080-gra"
  zone         = "all"
  port         = "10,11"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `display_name` - Human readable name for your frontend
* `port` - Port(s) attached to your frontend. Supports single port (numerical value), 
   range (2 dash-delimited increasing ports) and comma-separated list of 'single port' 
   and/or 'range'. Each port must be in the [1;49151] range
* `zone` - (Required) Zone where the frontend will be defined (ie. `gra`, `bhs` also supports `all`)
* `dedicated_ipfo` - Only attach frontend on these ip. No restriction if null. List of Ip blocks.
* `default_farm_id` - Default UDP Farm of your frontend
* `disabled` - Disable your frontend. Default: 'false'

## Attributes Reference

The following attributes are exported:

* `frontend_id` - Id of your frontend
* `display_name` - See Argument Reference above
* `port` - See Argument Reference above
* `zone` - See Argument Reference above
* `dedicated_ipfo` - See Argument Reference above
* `default_farm_id` - See Argument Reference above
* `disabled` - See Argument Reference above
