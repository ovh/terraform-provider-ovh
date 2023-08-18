---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing_tcp_frontend

Creates a backend server group (frontend) to be used by loadbalancing frontend(s)

## Example Usage

```hcl
data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_tcp_farm" "farm80" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name = "ingress-8080-gra"
  zone         = "all"
  port         = 80
}

resource "ovh_iploadbalancing_tcp_frontend" "testfrontend" {
  service_name    = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name    = "ingress-8080-gra"
  zone            = "all"
  port            = "80,443"
  default_farm_id = "${ovh_iploadbalancing_tcp_farm.farm80.id}"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `display_name` - Human readable name for your frontend, this field is for you
* `port` - Port(s) attached to your frontend. Supports single port (numerical value), 
   range (2 dash-delimited increasing ports) and comma-separated list of 'single port' 
   and/or 'range'. Each port must be in the [1;49151] range
* `zone` - (Required) Zone where the frontend will be defined (ie. `gra`, `bhs` also supports `all`)
* `allowed_source` - Restrict IP Load Balancing access to these ip block. No restriction if null. List of IP blocks.
* `dedicated_ipfo` - Only attach frontend on these ip. No restriction if null. List of Ip blocks.
* `default_farm_id` - Default TCP Farm of your frontend
* `default_ssl_id` - Default ssl served to your customer
* `disabled` - Disable your frontend. Default: 'false'
* `ssl` - SSL deciphering. Default: 'false'


## Attributes Reference

The following attributes are exported:

* `id` - Id of your frontend
* `display_name` - See Argument Reference above.
* `allowed_source` - See Argument Reference above.
* `dedicated_ipfo` - See Argument Reference above.
* `default_farm_id` - See Argument Reference above.
* `default_ssl_id` - See Argument Reference above.
* `disabled` - See Argument Reference above.
* `ssl` - See Argument Reference above.

## Import 

TCP frontend can be imported using the following format `service_name` and the `id` of the frontend separated by "/" e.g.

```bash
$ terraform import ovh_iploadbalancing_tcp_frontend.testfrontend service_name/tcp_frontend_id
```